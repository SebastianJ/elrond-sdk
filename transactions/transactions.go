package transactions

import (
	"encoding/hex"
	"math/big"
	"strings"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/hashing/blake2b"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/SebastianJ/elrond-sdk/api"
	"github.com/SebastianJ/elrond-sdk/utils"
	sdkWallet "github.com/SebastianJ/elrond-sdk/wallet"
)

var (
	TxJSONMarshaler     *marshal.TxJsonMarshalizer    = &marshal.TxJsonMarshalizer{}
	Hasher              *blake2b.Blake2b              = &blake2b.Blake2b{}
	InternalMarshalizer *marshal.GogoProtoMarshalizer = &marshal.GogoProtoMarshalizer{}
)

// Transaction - wrapper for transaction and API data
type Transaction struct {
	Transaction     transaction.Transaction
	APIData         api.TransactionData
	SenderShardID   uint32
	ReceiverShardID uint32
	TxHash          string
}

// SendTransaction - generates and broadcasts a transaction to the blockchain
func SendTransaction(
	wallet sdkWallet.Wallet,
	receiver string,
	amount float64,
	sendMaximumAmount bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (string, error) {
	tx, err := GenerateAndSignTransaction(wallet, receiver, amount, sendMaximumAmount, nonce, txData, gasParams, client)
	if err != nil {
		return "", err
	}

	txHexHash, txError := client.SendTransaction(tx.APIData)

	if txError != nil {
		// If we've sent an invalid nonce - sleep 3 seconds and then retry again using a fresh nonce
		if strings.Contains(txError.Error(), "transaction generation failed: invalid nonce") {
			time.Sleep(1 * time.Second)
			return SendTransaction(wallet, receiver, amount, sendMaximumAmount, -1, txData, gasParams, client)
		}

		return "", txError
	}

	return txHexHash, nil
}

// GenerateAndSignTransaction - generates and signs a transaction
func GenerateAndSignTransaction(
	wallet sdkWallet.Wallet,
	receiver string,
	amount float64,
	sendMaximumAmount bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (Transaction, error) {
	tx, err := GenerateTransaction(wallet, receiver, amount, sendMaximumAmount, nonce, txData, gasParams, client)

	signature, err := SignTransaction(wallet, tx)
	if err != nil {
		return Transaction{}, err
	}

	hexSignature := hex.EncodeToString(signature)
	tx.APIData.Signature = hexSignature

	return tx, nil
}

// GenerateTransaction - generates a new transaction using the supplied parameters
func GenerateTransaction(
	wallet sdkWallet.Wallet,
	receiver string,
	amount float64,
	sendMaximumAmount bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (Transaction, error) {
	receiverBytes, err := wallet.Converter.Decode(receiver)

	currentNonce, err := getNonce(client, wallet.Address, nonce)
	if err != nil {
		return Transaction{}, err
	}

	gasParams.UpdateGasLimit(txData)

	correctAmount, err := calculateAmount(client, wallet.Address, amount, sendMaximumAmount, gasParams)
	if err != nil {
		return Transaction{}, err
	}

	//converted, _ := utils.ConvertNumeralStringToBigFloat(realAmount.String())
	//fmt.Println(fmt.Sprintf("Sending amount: %f (%s)", converted, realAmount))

	innerTx := transaction.Transaction{
		SndAddr:  wallet.AddressBytes,
		RcvAddr:  receiverBytes,
		Value:    correctAmount,
		Data:     []byte(txData),
		Nonce:    currentNonce,
		GasPrice: gasParams.GasPrice,
		GasLimit: gasParams.GasLimit,
	}

	apiData := api.TransactionData{
		Sender:   wallet.Address,
		Receiver: receiver,
		Value:    correctAmount.String(),
		Data:     txData,
		Nonce:    currentNonce,
		GasPrice: gasParams.GasPrice,
		GasLimit: gasParams.GasLimit,
	}

	txHash, err := core.CalculateHash(InternalMarshalizer, Hasher, innerTx)
	if err != nil {
		return Transaction{}, err
	}
	txHexHash := hex.EncodeToString(txHash)

	tx := Transaction{
		Transaction: innerTx,
		APIData:     apiData,
		TxHash:      txHexHash,
	}

	return tx, nil
}

// SignTransaction - signs a given transaction and returns the signature
func SignTransaction(wallet sdkWallet.Wallet, tx Transaction) ([]byte, error) {
	txBuff, err := tx.Transaction.GetDataForSigning(wallet.Converter, TxJSONMarshaler)
	if err != nil {
		return nil, err
	}

	return wallet.Sign(txBuff)
}

func getNonce(client api.Client, address string, nonce int64) (currentNonce uint64, err error) {
	if nonce >= 0 {
		return uint64(nonce), nil
	}

	var account api.Account
	account, err = client.GetAccount(address)
	if err != nil {
		return 0, err
	}

	return uint64(account.Nonce), nil
}

func calculateAmount(client api.Client, address string, amount float64, sendMaximumAmount bool, gasParams GasParams) (correctAmount *big.Int, err error) {
	if !sendMaximumAmount {
		return utils.ConvertFloatAmountToBigInt(amount), nil
	}

	account, err := client.GetAccount(address)
	if err != nil {
		return nil, err
	}

	gasCost := gasParams.CalculateTotalGasCost()
	apiAmount, _ := new(big.Int).SetString(account.BalanceString, 10)
	correctAmount = gasParams.CalculateAmountWithoutGasCost(apiAmount, gasCost)

	return correctAmount, nil
}

package transactions

import (
	"encoding/hex"
	"math/big"
	"strings"
	"time"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/SebastianJ/elrond-sdk/api"
	"github.com/SebastianJ/elrond-sdk/utils"
	sdkWallet "github.com/SebastianJ/elrond-sdk/wallet"
)

// SendTransaction - generates and broadcasts a transaction to the blockchain
func SendTransaction(
	wallet sdkWallet.Wallet,
	receiver string,
	amount float64,
	maximum bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (string, error) {
	_, apiData, err := GenerateAndSignTransaction(wallet, receiver, amount, maximum, nonce, txData, gasParams, client)
	if err != nil {
		return "", err
	}

	txHexHash, txError := client.SendTransaction(apiData)

	if txError != nil {
		// If we've sent an invalid nonce - sleep 3 seconds and then retry again using a fresh nonce
		if strings.Contains(txError.Error(), "transaction generation failed: invalid nonce") {
			time.Sleep(3 * time.Second)
			return SendTransaction(wallet, receiver, amount, maximum, nonce, txData, gasParams, client)
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
	maximum bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (transaction.Transaction, api.TransactionData, error) {
	tx, apiData, err := GenerateTransaction(wallet, receiver, amount, maximum, nonce, txData, gasParams, client)

	signature, err := signTransaction(wallet, tx)
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	hexSignature := hex.EncodeToString(signature)
	apiData.Signature = hexSignature

	return tx, apiData, nil
}

// GenerateTransaction - generates a new transaction using the supplied parameters
func GenerateTransaction(
	wallet sdkWallet.Wallet,
	receiver string,
	amount float64,
	maximum bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (transaction.Transaction, api.TransactionData, error) {
	senderBytes, err := wallet.PrivateKey.GeneratePublic().ToByteArray()
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	sender := wallet.Converter.Encode(senderBytes)
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	receiverBytes, err := wallet.Converter.Decode(receiver)

	account, err := client.GetAccount(sender)
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	realNonce := getNonce(&account, nonce)

	if len(txData) > 0 {
		gasParams.GasLimit = gasParams.GasLimit + (uint64(len(txData)) * gasParams.GasPerDataByte)
	}

	var realAmount *big.Int
	if maximum {
		realAmount = calculateMaximumAmount(&account, gasParams.GasPrice, gasParams.GasLimit)
	} else {
		realAmount = utils.ConvertFloatAmountToBigInt(amount)
	}

	//converted, _ := utils.ConvertNumeralStringToBigFloat(realAmount.String())
	//fmt.Println(fmt.Sprintf("Sending amount: %f (%s)", converted, realAmount))

	tx := transaction.Transaction{
		SndAddr:  senderBytes,
		RcvAddr:  receiverBytes,
		Value:    realAmount,
		Data:     []byte(txData),
		Nonce:    realNonce,
		GasPrice: gasParams.GasPrice,
		GasLimit: gasParams.GasLimit,
	}

	apiData := api.TransactionData{
		Sender:   sender,
		Receiver: receiver,
		Value:    realAmount.String(),
		Data:     txData,
		Nonce:    realNonce,
		GasPrice: gasParams.GasPrice,
		GasLimit: gasParams.GasLimit,
	}

	return tx, apiData, nil
}

func signTransaction(wallet sdkWallet.Wallet, tx transaction.Transaction) ([]byte, error) {
	marshaler := &marshal.TxJsonMarshalizer{}
	txBuff, err := tx.GetDataForSigning(wallet.Converter, marshaler)
	if err != nil {
		return nil, err
	}

	return wallet.Sign(txBuff)
}

func getNonce(accountData *api.Account, nonce int64) uint64 {
	var realNonce uint64

	if nonce > 0 {
		realNonce = uint64(nonce)
	} else if accountData != nil {
		realNonce = accountData.Nonce
	}

	return realNonce
}

func calculateMaximumAmount(accountData *api.Account, gasPrice uint64, gasLimit uint64) *big.Int {
	gasCost := utils.CalculateTotalGasCost(gasPrice, gasLimit)
	apiAmount, _ := new(big.Int).SetString(accountData.Balance, 10)
	realAmount := utils.CalculateAmountWithoutGasCost(apiAmount, gasCost)

	return realAmount
}

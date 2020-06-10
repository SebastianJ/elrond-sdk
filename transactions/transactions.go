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
	sendMaximumAmount bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (string, error) {
	_, apiData, err := GenerateAndSignTransaction(wallet, receiver, amount, sendMaximumAmount, nonce, txData, gasParams, client)
	if err != nil {
		return "", err
	}

	txHexHash, txError := client.SendTransaction(apiData)

	if txError != nil {
		// If we've sent an invalid nonce - sleep 3 seconds and then retry again using a fresh nonce
		if strings.Contains(txError.Error(), "transaction generation failed: invalid nonce") {
			time.Sleep(3 * time.Second)
			return SendTransaction(wallet, receiver, amount, sendMaximumAmount, nonce, txData, gasParams, client)
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
) (transaction.Transaction, api.TransactionData, error) {
	tx, apiData, err := GenerateTransaction(wallet, receiver, amount, sendMaximumAmount, nonce, txData, gasParams, client)

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
	sendMaximumAmount bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (transaction.Transaction, api.TransactionData, error) {
	receiverBytes, err := wallet.Converter.Decode(receiver)

	currentNonce, err := getNonce(client, wallet.Address, nonce)
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	gasParams.UpdateGasLimit(txData)

	correctAmount, err := calculateAmount(client, wallet.Address, amount, sendMaximumAmount, gasParams)
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	//converted, _ := utils.ConvertNumeralStringToBigFloat(realAmount.String())
	//fmt.Println(fmt.Sprintf("Sending amount: %f (%s)", converted, realAmount))

	tx := transaction.Transaction{
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

func getNonce(client api.Client, address string, nonce int64) (currentNonce uint64, err error) {
	var account api.Account

	if nonce > 0 {
		currentNonce = uint64(nonce)
	} else {
		account, err = client.GetAccount(address)
		if err != nil {
			return 0, err
		}
		currentNonce = uint64(account.Nonce)
	}

	return currentNonce, err
}

func calculateAmount(client api.Client, address string, amount float64, sendMaximumAmount bool, gasParams GasParams) (correctAmount *big.Int, err error) {
	if sendMaximumAmount {
		account, err := client.GetAccount(address)
		if err != nil {
			return nil, err
		}

		gasCost := gasParams.CalculateTotalGasCost()
		apiAmount, _ := new(big.Int).SetString(account.Balance, 10)
		correctAmount = gasParams.CalculateAmountWithoutGasCost(apiAmount, gasCost)
	} else {
		correctAmount = utils.ConvertFloatAmountToBigInt(amount)
	}

	return correctAmount, nil
}

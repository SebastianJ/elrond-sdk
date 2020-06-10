package transactions

import (
	"encoding/hex"
	"math/big"
	"strings"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519/singlesig"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/SebastianJ/elrond-sdk/api"
	cliCrypto "github.com/SebastianJ/elrond-sdk/crypto"
	"github.com/SebastianJ/elrond-sdk/utils"
)

var (
	converter core.PubkeyConverter
)

func init() {
	converter, _ = pubkeyConverter.NewBech32PubkeyConverter(32)
}

// SendTransaction - generates and broadcasts a transaction to the blockchain
func SendTransaction(
	walletPath string,
	receiver string,
	amount float64,
	maximum bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (string, error) {
	signer, privateKey, _, err := cliCrypto.DecryptWallet(walletPath)
	if err != nil {
		return "", err
	}

	_, apiData, err := GenerateAndSignTransaction(privateKey, signer, receiver, amount, maximum, nonce, txData, gasParams, client)
	client.Initialize()
	txHexHash, txError := client.SendTransaction(apiData)

	if txError != nil {
		// If we've sent an invalid nonce - sleep 3 seconds and then retry again using a fresh nonce
		if strings.Contains(txError.Error(), "transaction generation failed: invalid nonce") {
			time.Sleep(3 * time.Second)
			return SendTransaction(walletPath, receiver, amount, maximum, nonce, txData, gasParams, client)
		}

		return "", txError
	}

	return txHexHash, nil
}

// GenerateAndSignTransaction - generates and signs a transaction
func GenerateAndSignTransaction(
	privateKey crypto.PrivateKey,
	signer *singlesig.Ed25519Signer,
	receiver string,
	amount float64,
	maximum bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (transaction.Transaction, api.TransactionData, error) {
	tx, apiData, err := GenerateTransaction(privateKey, receiver, amount, maximum, nonce, txData, gasParams, client)

	signature, err := signTransaction(privateKey, signer, tx)
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	hexSignature := hex.EncodeToString(signature)
	apiData.Signature = hexSignature

	return tx, apiData, nil
}

// GenerateTransaction - generates a new transaction using the supplied parameters
func GenerateTransaction(
	privateKey crypto.PrivateKey,
	receiver string,
	amount float64,
	maximum bool,
	nonce int64,
	txData string,
	gasParams GasParams,
	client api.Client,
) (transaction.Transaction, api.TransactionData, error) {
	senderBytes, err := privateKey.GeneratePublic().ToByteArray()
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	sender := converter.Encode(senderBytes)
	if err != nil {
		return transaction.Transaction{}, api.TransactionData{}, err
	}

	receiverBytes, err := converter.Decode(receiver)

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

func signTransaction(privateKey crypto.PrivateKey, signer *singlesig.Ed25519Signer, tx transaction.Transaction) ([]byte, error) {
	marshaler := &marshal.TxJsonMarshalizer{}
	txBuff, err := tx.GetDataForSigning(converter, marshaler)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(privateKey, txBuff)
	if err != nil {
		return nil, err
	}

	return signature, nil
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

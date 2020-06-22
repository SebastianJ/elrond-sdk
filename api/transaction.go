package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// TransactionData - represents the transaction data sent to a node API to send transactions
type TransactionData struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Value     string `json:"value"`
	Data      string `json:"data"`
	Nonce     uint64 `json:"nonce"`
	GasPrice  uint64 `json:"gasPrice"`
	GasLimit  uint64 `json:"gasLimit"`
	Signature string `json:"signature"`
}

// SendTransactionResponse - API response when sending one transaction
type SendTransactionResponse struct {
	TxHash string `json:"txHash"`
	Error  string `json:"error,omitempty"`
}

// SendMultipleTransactionsResponse - API response when sending multiple transactions
type SendMultipleTransactionsResponse struct {
	TxsSent   uint64         `json:"txsSent"`
	TxsHashes map[int]string `json:"txsHashes"`
	Error     string         `json:"error,omitempty"`
}

// SendTransaction performs the actual HTTP request to send the transaction
func (client *Client) SendTransaction(txData *TransactionData) (string, error) {
	client.Initialize()

	url := fmt.Sprintf("%s/transaction/send", client.Host)

	jsonData, err := json.Marshal(txData)
	if err != nil {
		return "", errors.Wrapf(err, "JSON Marshal")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.Wrapf(err, "HTTP NewRequest")
	}

	body, err := client.PerformRequest(url, req)
	if err != nil {
		return "", errors.Wrapf(err, "Client PerformRequest")
	}

	var response SendTransactionResponse
	json.Unmarshal([]byte(body), &response)

	if response.Error != "" {
		return "", fmt.Errorf("Response error: %s", response.Error)
	}

	return response.TxHash, nil
}

// SendMultipleTransactions performs the actual HTTP request to send the transactions
func (client *Client) SendMultipleTransactions(txs []*TransactionData) (SendMultipleTransactionsResponse, error) {
	client.Initialize()

	url := fmt.Sprintf("%s/transaction/send-multiple", client.Host)

	jsonData, err := json.Marshal(txs)
	if err != nil {
		return SendMultipleTransactionsResponse{}, errors.Wrapf(err, "JSON Marshal")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return SendMultipleTransactionsResponse{}, errors.Wrapf(err, "HTTP NewRequest")
	}

	body, err := client.PerformRequest(url, req)
	if err != nil {
		return SendMultipleTransactionsResponse{}, errors.Wrapf(err, "Client PerformRequest")
	}

	var response SendMultipleTransactionsResponse
	json.Unmarshal([]byte(body), &response)

	if response.Error != "" {
		return SendMultipleTransactionsResponse{}, fmt.Errorf("Response error: %s", response.Error)
	}

	return response, nil
}

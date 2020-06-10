package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

type sendTxResponse struct {
	TxHash string `json:"txHash"`
	Error  string `json:"error,omitempty"`
}

// SendTransaction performs the actual HTTP request to send the transaction
func (client *Client) SendTransaction(txData TransactionData) (string, error) {
	url := fmt.Sprintf("%s/transaction/send", client.Host)

	jsonData, err := json.Marshal(txData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	body, err := client.PerformRequest(url, req)
	if err != nil {
		return "", err
	}

	var response sendTxResponse
	json.Unmarshal([]byte(body), &response)

	if response.TxHash == "" {
		return "", errors.New(response.Error)
	}

	return response.TxHash, nil
}

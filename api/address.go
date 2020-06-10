package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Account contains the current data for a specific wallet or account
type Account struct {
	Address string `json:"address,omitempty"`
	Nonce   uint64 `json:"nonce,omitempty"`
	Balance string `json:"balance"`
	//Code     string
	//CodeHash []byte
	//RootHash []byte
}

// AccountWrapper is simple wrapper type to help with deserializing the address response
type AccountWrapper struct {
	Account Account `json:"account"`
}

// GetAccount fetches the desired account's balance as well as nonce
func (client *Client) GetAccount(address string) (Account, error) {
	host := client.Host

	if client.ForceAPINonceLookups {
		host = defaultEndpoint
	}

	url := fmt.Sprintf("%s/address/%s", host, address)
	req, err := http.NewRequest("GET", url, nil)

	var response AccountWrapper
	var accountResponse Account

	body, err := client.PerformRequest(url, req)

	if err != nil {
		return accountResponse, err
	}

	json.Unmarshal([]byte(body), &response)
	accountResponse = response.Account

	return accountResponse, nil
}

// GetBalance fetches the balance of a specific account
func (client *Client) GetBalance(address string) (Account, error) {
	host := client.Host

	if client.ForceAPINonceLookups {
		host = defaultEndpoint
	}

	url := fmt.Sprintf("%s/address/%s/balance", host, address)
	req, err := http.NewRequest("GET", url, nil)

	var accountResponse Account

	if err != nil {
		return accountResponse, err
	}

	body, err := client.PerformRequest(url, req)

	if err != nil {
		return accountResponse, err
	}

	json.Unmarshal([]byte(body), &accountResponse)

	return accountResponse, nil
}

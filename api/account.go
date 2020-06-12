package api

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/SebastianJ/elrond-sdk/utils"
)

// Account contains the current data for a specific wallet or account
type Account struct {
	Address       string     `json:"address,omitempty"`
	Nonce         uint64     `json:"nonce,omitempty"`
	BalanceString string     `json:"balance"`
	Balance       *big.Float `json:"-"`
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
	client.Initialize()

	host := client.Host

	if client.ForceAPINonceLookups {
		host = defaultEndpoint
	}

	url := fmt.Sprintf("%s/address/%s", host, address)
	req, err := http.NewRequest("GET", url, nil)

	var response AccountWrapper
	var account Account

	body, err := client.PerformRequest(url, req)

	if err != nil {
		return account, err
	}

	json.Unmarshal([]byte(body), &response)
	account = response.Account

	if err := account.Initialize(); err != nil {
		return account, err
	}

	return account, nil
}

// GetBalance fetches the balance of a specific account
func (client *Client) GetBalance(address string) (Account, error) {
	client.Initialize()

	host := client.Host

	if client.ForceAPINonceLookups {
		host = defaultEndpoint
	}

	var account Account
	url := fmt.Sprintf("%s/address/%s/balance", host, address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return account, err
	}

	body, err := client.PerformRequest(url, req)

	if err != nil {
		return account, err
	}

	json.Unmarshal([]byte(body), &account)

	if err := account.Initialize(); err != nil {
		return account, err
	}

	return account, nil
}

// Initialize - convert balances etc
func (account *Account) Initialize() error {
	if account.BalanceString != "" {
		converted, err := utils.ConvertNumeralStringToBigFloat(account.BalanceString)
		if err != nil {
			return err
		}

		account.Balance = converted
	}

	return nil
}

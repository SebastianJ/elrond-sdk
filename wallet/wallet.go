package wallet

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/SebastianJ/elrond-sdk/crypto"
)

var (
	bech32KeyLength = 32
)

// Wallet - represents a wallet
type Wallet struct {
	Key          crypto.Key
	Converter    core.PubkeyConverter
	Address      string
	AddressBytes []byte
}

// Decrypt - decrypts a given PEM wallet file
func Decrypt(walletPath string) (Wallet, error) {
	encodedSk, _, err := core.LoadSkPkFromPemFile(walletPath, 0)
	if err != nil {
		return Wallet{}, err
	}

	skBytes, err := hex.DecodeString(string(encodedSk))
	if err != nil {
		return Wallet{}, fmt.Errorf("%w for encoded secret key", err)
	}

	key, err := crypto.LoadKeyFromPrivateKey(crypto.ED25519, skBytes)
	if err != nil {
		return Wallet{}, err
	}

	return toWallet(key)
}

// Generate - generate a new wallet
func Generate() (Wallet, error) {
	key, err := crypto.GenerateKey(crypto.ED25519)
	if err != nil {
		return Wallet{}, err
	}

	return toWallet(key)
}

func toWallet(key crypto.Key) (Wallet, error) {
	converter, err := pubkeyConverter.NewBech32PubkeyConverter(bech32KeyLength)
	if err != nil {
		return Wallet{}, err
	}

	addressBytes, err := key.PrivateKey.GeneratePublic().ToByteArray()
	if err != nil {
		return Wallet{}, err
	}

	address := converter.Encode(addressBytes)
	if err != nil {
		return Wallet{}, err
	}

	wallet := Wallet{
		Key:          key,
		Converter:    converter,
		Address:      address,
		AddressBytes: addressBytes,
	}

	return wallet, nil
}

// Sign - sign provided []byte data using the wallet
func (wallet *Wallet) Sign(data []byte) ([]byte, error) {
	signature, err := wallet.Key.Signer.Sign(wallet.Key.PrivateKey, data)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

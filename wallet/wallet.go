package wallet

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519/singlesig"
)

// GasParams - represents gas parameters for a transaction
type Wallet struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
	Signer     *singlesig.Ed25519Signer
	Converter  core.PubkeyConverter
}

func Decrypt(walletPath string) (Wallet, error) {
	encodedSk, _, err := core.LoadSkPkFromPemFile(walletPath, 0)
	if err != nil {
		return Wallet{}, err
	}

	skBytes, err := hex.DecodeString(string(encodedSk))
	if err != nil {
		return Wallet{}, fmt.Errorf("%w for encoded secret key", err)
	}

	signer, privKey, pubKey, err := generateCryptoSuite(skBytes)
	if err != nil {
		return Wallet{}, err
	}

	converter, err := pubkeyConverter.NewBech32PubkeyConverter(32)
	if err != nil {
		return Wallet{}, err
	}

	wallet := Wallet{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		Signer:     signer,
		Converter:  converter,
	}

	return wallet, nil
}

func generateCryptoSuite(skBytes []byte) (signer *singlesig.Ed25519Signer, privKey crypto.PrivateKey, pubKey crypto.PublicKey, err error) {
	signer = &singlesig.Ed25519Signer{}
	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())

	privKey, err = keyGen.PrivateKeyFromByteArray(skBytes)
	if err != nil {
		return nil, nil, nil, err
	}

	pubKey = privKey.GeneratePublic()

	return signer, privKey, pubKey, err
}

// Sign - sign provided []byte data using the wallet
func (wallet *Wallet) Sign(data []byte) ([]byte, error) {
	signature, err := wallet.Signer.Sign(wallet.PrivateKey, data)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

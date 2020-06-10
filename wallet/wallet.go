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

// Wallet - represents a wallet
type Wallet struct {
	PrivateKey   crypto.PrivateKey
	PublicKey    crypto.PublicKey
	Signer       *singlesig.Ed25519Signer
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

	return toWallet(skBytes)
}

// Generate - generate a new wallet
func Generate() (Wallet, error) {
	keyGen := newKeyGenerator()

	sk, _ := keyGen.GeneratePair()
	skBytes, err := sk.ToByteArray()
	if err != nil {
		return Wallet{}, fmt.Errorf("%w for encoded secret key", err)
	}

	return toWallet(skBytes)
}

func toWallet(skBytes []byte) (Wallet, error) {
	signer, privateKey, publicKey, err := generateCryptoSuite(skBytes)
	if err != nil {
		return Wallet{}, err
	}

	converter, err := pubkeyConverter.NewBech32PubkeyConverter(32)
	if err != nil {
		return Wallet{}, err
	}

	addressBytes, err := privateKey.GeneratePublic().ToByteArray()
	if err != nil {
		return Wallet{}, err
	}

	address := converter.Encode(addressBytes)
	if err != nil {
		return Wallet{}, err
	}

	wallet := Wallet{
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		Signer:       signer,
		Converter:    converter,
		Address:      address,
		AddressBytes: addressBytes,
	}

	return wallet, nil
}

func newSigner() *singlesig.Ed25519Signer {
	return &singlesig.Ed25519Signer{}
}

func newKeyGenerator() crypto.KeyGenerator {
	return signing.NewKeyGenerator(ed25519.NewEd25519())
}

func generateCryptoSuite(skBytes []byte) (signer *singlesig.Ed25519Signer, privKey crypto.PrivateKey, pubKey crypto.PublicKey, err error) {
	signer = newSigner()
	keyGen := newKeyGenerator()

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

package crypto

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519/singlesig"
)

func DecryptWallet(walletPath string) (*singlesig.Ed25519Signer, crypto.PrivateKey, crypto.PublicKey, error) {
	encodedSk, _, err := core.LoadSkPkFromPemFile(walletPath, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	skBytes, err := hex.DecodeString(string(encodedSk))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w for encoded secret key", err)
	}

	signer, privKey, pubKey, err := generateCryptoSuite(skBytes)
	if err != nil {
		return nil, nil, nil, err
	}

	return signer, privKey, pubKey, nil
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

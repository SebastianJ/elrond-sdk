package crypto

import (
	"github.com/ElrondNetwork/elrond-go/core"
	erdCrypto "github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519"
	erdSingleSig "github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519/singlesig"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/mcl"
	blsSinglesig "github.com/ElrondNetwork/elrond-go/crypto/signing/mcl/singlesig"
)

const (
	// BLS12 cipher
	BLS12 = iota
	// ED25519 cipher
	ED25519
)

// Key - represents a cryptographic key
type Key struct {
	Cipher           int
	PrivateKey       erdCrypto.PrivateKey
	PublicKey        erdCrypto.PublicKey
	PrivateKeyBytes  []byte
	PrivateKeyString string
	PublicKeyBytes   []byte
	PublicKeyString  string
	Signer           erdCrypto.SingleSigner
	Converter        core.PubkeyConverter
}

// LoadKeyFromPrivateKey - instantiate a new key based on supplied private key bytes
func LoadKeyFromPrivateKey(cipher int, privateKeyBytes []byte) (Key, error) {
	keyGen := NewKeyGenerator(cipher)

	privateKey, err := keyGen.PrivateKeyFromByteArray(privateKeyBytes)
	if err != nil {
		return Key{}, err
	}

	publicKey := privateKey.GeneratePublic()

	return NewKey(cipher, privateKey, publicKey)
}

// GenerateKey - generate a new key pair based on a supplied cipher
func GenerateKey(cipher int) (Key, error) {
	keyGen := NewKeyGenerator(cipher)
	privateKey, publicKey := keyGen.GeneratePair()

	return NewKey(cipher, privateKey, publicKey)
}

// NewKey - Generate a new key instance based on a supplied cipher, private key and public key
func NewKey(cipher int, privateKey erdCrypto.PrivateKey, publicKey erdCrypto.PublicKey) (Key, error) {
	key := Key{
		Cipher:     cipher,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Signer:     NewSigner(cipher),
	}

	if err := key.SetKeyBytes(); err != nil {
		return key, err
	}

	return key, nil
}

// NewSigner - generate a new signer based on supplied cipher
func NewSigner(cipher int) erdCrypto.SingleSigner {
	if cipher == BLS12 {
		return blsSinglesig.NewBlsSigner()
	}

	return &erdSingleSig.Ed25519Signer{}
}

// NewKeyGenerator - generate a new key generator based on supplied cipher
func NewKeyGenerator(cipher int) erdCrypto.KeyGenerator {
	if cipher == BLS12 {
		return signing.NewKeyGenerator(mcl.NewSuiteBLS12())
	}

	return signing.NewKeyGenerator(ed25519.NewEd25519())
}

// SetKeyBytes - converts the key's private and public keys to byte arrays
func (key *Key) SetKeyBytes() (err error) {
	key.PrivateKeyBytes, err = keyToBytes(key.PrivateKey)
	if err != nil {
		return err
	}

	key.PublicKeyBytes, err = keyToBytes(key.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

func keyToBytes(key erdCrypto.Key) ([]byte, error) {
	bytes, err := key.ToByteArray()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

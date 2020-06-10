package keystore

import (
	"io/ioutil"
	"os"

	ethKeystore "github.com/ethereum/go-ethereum/accounts/keystore"
)

func DecryptKeystoreFile(path string, auth string) (*ethKeystore.Key, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	jsonData, _ := ioutil.ReadAll(file)

	return ethKeystore.DecryptKey(jsonData, auth)
}

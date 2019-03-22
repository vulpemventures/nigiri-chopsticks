package regtestfaucet

import (
	"github.com/vulpemventures/nigiri-chopsticks/faucet"
)

type regtestfaucet struct {
	URL string
}

func (f *regtestfaucet) New(url string) {
	f.URL = url
}

func (f *regtestfaucet) Send(address string) (int, string, error) {
	status, utxo, err := listunspent(f.URL)
	if err != nil {
		return status, "", err
	}

	status, tx, err := createrawtransaction(f.URL, utxo, address)
	if err != nil {
		return status, "", err
	}

	status, signedTx, err := signrawtransaction(f.URL, tx)
	if err != nil {
		return status, "", err
	}

	return status, signedTx, nil
}

// NewFaucet initialize a regtest faucet and returns it as interface
func NewFaucet(url string) faucet.Faucet {
	f := &regtestfaucet{}
	f.New(url)

	return f
}

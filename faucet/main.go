package faucet

type Faucet interface {
	NewTransaction(address string) (int, string, error)
}

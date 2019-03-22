package faucet

type Faucet interface {
	Send(address string) (int, string, error)
}

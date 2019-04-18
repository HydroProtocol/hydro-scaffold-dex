package cli

//TODO move to backend sdk
type Erc20 interface {
	Symbol(address string) (error, string)
	Decimals(address string) (error, int)
	Name(address string) (error, string)
}

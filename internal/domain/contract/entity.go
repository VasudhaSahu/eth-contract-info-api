package contract

type Address string

type Info struct {
	Address      Address
	Network      string
	BlockTag     string
	HasCode      bool
	Bytecode     string
	BytecodeSize int
}

type ERC20Metadata struct {
	Address     Address
	Network     string
	BlockTag    string
	Name        string
	Symbol      string
	Decimals    int
	TotalSupply string
}

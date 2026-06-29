package contract

type GetContractInfoInput struct {
	Address  string
	BlockTag string
}

type GetERC20MetadataInput struct {
	Address  string
	BlockTag string
}

type ContractInfoResponse struct {
	Address         string `json:"address"`
	Network         string `json:"network"`
	BlockTag        string `json:"blockTag"`
	HasCode         bool   `json:"hasCode"`
	BytecodeSize    int    `json:"bytecodeSize"`
	BytecodePreview string `json:"bytecodePreview"`
}

type ERC20MetadataResponse struct {
	Address     string `json:"address"`
	Network     string `json:"network"`
	BlockTag    string `json:"blockTag"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    int    `json:"decimals"`
	TotalSupply string `json:"totalSupply"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

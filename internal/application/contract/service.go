package contract

import d "github.com/VasudhaSahu/eth-contract-info-api.git/internal/domain/contract"

type Reader interface {
	GetCode(address string, blockTag string) (string, error)
	Call(address string, data string, blockTag string) (string, error)
}

type UseCase struct {
	Reader  Reader
	Network string
}

func NewUseCase(reader Reader, network string) *UseCase {
	return &UseCase{Reader: reader, Network: network}
}

// GetInfo is intentionally narrow: one RPC call, a few domain rules, and a DTO.
func (u *UseCase) GetInfo(input GetContractInfoInput) (ContractInfoResponse, error) {

	code, err := u.Reader.GetCode(input.Address, input.BlockTag)
	if err != nil {
		return ContractInfoResponse{}, err
	}

	return ContractInfoResponse{
		Address:         string(d.NormalizeAddress(input.Address)),
		Network:         u.Network,
		BlockTag:        input.BlockTag,
		HasCode:         d.HasContractCode(code),
		BytecodeSize:    d.BytecodeSize(code),
		BytecodePreview: d.BytecodePreview(code, 66),
	}, nil
}

func (u *UseCase) GetERC20Metadata(input GetERC20MetadataInput) (ERC20MetadataResponse, error) {

	// ERC‑20 metadata is encoded as a series of eth_call responses, we decode only the basics here.
	nameRaw, err := u.Reader.Call(input.Address, d.SelectorName, input.BlockTag)
	if err != nil {
		return ERC20MetadataResponse{}, err
	}
	symbolRaw, err := u.Reader.Call(input.Address, d.SelectorSymbol, input.BlockTag)
	if err != nil {
		return ERC20MetadataResponse{}, err
	}
	decimalsRaw, err := u.Reader.Call(input.Address, d.SelectorDecimals, input.BlockTag)
	if err != nil {
		return ERC20MetadataResponse{}, err
	}
	totalSupplyRaw, err := u.Reader.Call(input.Address, d.SelectorTotalSupply, input.BlockTag)
	if err != nil {
		return ERC20MetadataResponse{}, err
	}

	return ERC20MetadataResponse{
		Address:     string(d.NormalizeAddress(input.Address)),
		Network:     u.Network,
		BlockTag:    input.BlockTag,
		Name:        decodeString(nameRaw),
		Symbol:      decodeString(symbolRaw),
		Decimals:    decodeUint(decimalsRaw),
		TotalSupply: decodeBigIntString(totalSupplyRaw),
	}, nil
}

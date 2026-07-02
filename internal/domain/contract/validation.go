package contract

import "strings"

const (
	SelectorName        = "0x06fdde03"
	SelectorSymbol      = "0x95d89b41"
	SelectorDecimals    = "0x313ce567"
	SelectorTotalSupply = "0x18160ddd"
)

func NormalizeAddress(address string) Address {
	return Address(strings.ToLower(strings.TrimSpace(address)))
}

func HasContractCode(code string) bool {
	return code != "" && code != "0x" && code != "0x0"
}

func BytecodeSize(code string) int {
	if !HasContractCode(code) {
		return 0
	}
	return len(strings.TrimPrefix(code, "0x")) / 2
}

func BytecodePreview(code string, max int) string {
	if len(code) <= max {
		return code
	}
	return code[:max] + "..."
}

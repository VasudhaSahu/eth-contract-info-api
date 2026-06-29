package contract

import (
	"encoding/hex"
	"math/big"
	"strings"
)

// decodeString expects a single dynamic string return value and uses the standard ABI layout.
// TODO: This is intentionally minimal and does not handle nested tuples or arrays.
func decodeString(data string) string {
	hexData := strings.TrimPrefix(data, "0x")
	if len(hexData) < 128 {
		return ""
	}
	if len(hexData) >= 128 {
		lengthHex := hexData[64:128]
		length := new(big.Int)
		length.SetString(lengthHex, 16)
		strLen := int(length.Int64())
		start := 128
		end := start + strLen*2
		if end > len(hexData) || strLen < 0 {
			return ""
		}
		b, err := hex.DecodeString(hexData[start:end])
		if err != nil {
			return ""
		}
		return string(b)
	}
	return ""
}

// decodeUint expects a single unsigned integer return value and uses the standard ABI layout.
func decodeUint(data string) int {
	hexData := strings.TrimPrefix(data, "0x")
	if hexData == "" {
		return 0
	}
	n := new(big.Int)
	n.SetString(hexData, 16)
	return int(n.Int64())
}

func decodeBigIntString(data string) string {
	hexData := strings.TrimPrefix(data, "0x")
	if hexData == "" {
		return "0"
	}
	n := new(big.Int)
	n.SetString(hexData, 16)
	return n.String()
}

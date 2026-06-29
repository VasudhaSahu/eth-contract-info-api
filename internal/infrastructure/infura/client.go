package infura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Client  *http.Client
}

type rpcRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  string    `json:"result"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// newClient creates a new Infura client with the given network and API key.
// NewClient is thin on purpose, the demo does not try to auto‑discover networks or endpoints.
func NewClient(network, apiKey string) *Client {
	// Basic sanity check; we’d fail fast if the API key is missing.
	if apiKey == "" {
		// For a real service we’d return an error instead of panicking.
		panic("INFURA_API_KEY is empty, check your configuration")
	}

	return &Client{
		BaseURL: fmt.Sprintf("https://%s.infura.io/v3/%s", network, apiKey),
		Client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// GetCode implements Reader.GetCode and delegates to doRPC.
// eth_getCode only tells you whether code is deployed and returns bytecode; it does not decode ABI methods, token metadata, or contract source automatically.
func (c *Client) GetCode(address string, blockTag string) (string, error) {
	return c.doRPC("eth_getCode", []interface{}{address, blockTag})
}

// Call implements Reader.Call and delegates to doRPC.
// eth_call can read token fields like name, symbol, decimals, and totalSupply without sending a transaction
func (c *Client) Call(address string, data string, blockTag string) (string, error) {
	callObj := map[string]string{
		"to":   address,
		"data": data,
	}
	return c.doRPC("eth_call", []interface{}{callObj, blockTag})
}

// doRPC performs a JSON-RPC call to the Infura API.
func (c *Client) doRPC(method string, params interface{}) (string, error) {

	// For this demo we always send a single JSON‑RPC request per HTTP call.
	// TODO: A batch API would be more efficient for bulk metadata fetches.
	reqBody := rpcRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal rpc request for %s: %w", method, err)
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("build http request for %s: %w", method, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "eth-contract-info-api/1.0")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call infura for %s: %w", method, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read infura response for %s: %w", method, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"infura returned status %d for %s: %s",
			resp.StatusCode, method, string(respBytes),
		)
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(respBytes, &rpcResp); err != nil {
		return "", fmt.Errorf("unmarshal rpc response for %s: %w", method, err)
	}

	if rpcResp.Error != nil {
		return "", fmt.Errorf(
			"rpc error for %s: code=%d msg=%s",
			method, rpcResp.Error.Code, rpcResp.Error.Message,
		)
	}

	return rpcResp.Result, nil
}

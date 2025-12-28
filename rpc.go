package chert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RPCClient handles JSON-RPC communication with the blockchain
type RPCClient struct {
	endpoint string
	client   *http.Client
}

// NewRPCClient creates a new RPC client
func NewRPCClient(endpoint string, timeout time.Duration) *RPCClient {
	return &RPCClient{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Call makes a JSON-RPC call to the blockchain
func (c *RPCClient) Call(ctx context.Context, method string, params interface{}, result interface{}) error {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal RPC request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create RPC request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("RPC request failed: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("failed to decode RPC response: %w", err)
	}

	if rpcResp.Error != nil {
		return rpcResp.Error
	}

	if result != nil && rpcResp.Result != nil {
		resultBytes, err := json.Marshal(rpcResp.Result)
		if err != nil {
			return fmt.Errorf("failed to marshal RPC result: %w", err)
		}
		return json.Unmarshal(resultBytes, result)
	}

	return nil
}
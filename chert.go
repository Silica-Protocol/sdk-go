// Package chert provides the official Go SDK for the Chert/Silica blockchain network.
//
// This SDK provides a comprehensive interface for interacting with the Chert blockchain,
// including wallet management, privacy features, staking, and governance operations.
//
// Basic usage:
//
//	import "github.com/silica-network/chert/sdk/go"
//
//	client, err := chert.NewClient(&chert.ClientConfig{
//		Endpoint: "https://api.chert.com",
//		Network:  chert.NetworkMainnet,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	account, err := client.Wallet.CreateAccount()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Created account: %s\n", account.Address)
package chert

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	// SDKVersion represents the current SDK version
	SDKVersion = "0.1.0"

	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second

	// DefaultEndpoint is the default API endpoint
	DefaultEndpoint = "https://api.chert.com"
)

// Network represents the blockchain network
type Network string

const (
	NetworkMainnet Network = "mainnet"
	NetworkTestnet Network = "testnet"
	NetworkDevnet  Network = "devnet"
)

// ClientConfig holds the configuration for the Chert client
type ClientConfig struct {
	// Endpoint is the API endpoint URL
	Endpoint string `json:"endpoint"`

	// Network specifies the blockchain network
	Network Network `json:"network"`

	// Timeout is the request timeout duration
	Timeout time.Duration `json:"timeout"`

	// APIKey is the optional API key for authenticated requests
	APIKey string `json:"api_key,omitempty"`

	// Headers contains additional HTTP headers
	Headers map[string]string `json:"headers,omitempty"`
}

// DefaultClientConfig returns a default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Endpoint: DefaultEndpoint,
		Network:  NetworkMainnet,
		Timeout:  DefaultTimeout,
		Headers:  make(map[string]string),
	}
}

// ChertClient is the main client for interacting with the Chert blockchain
type ChertClient struct {
	config     *ClientConfig
	httpClient *http.Client
	rpcClient  *RPCClient

	// Managers
	Wallet    *WalletManager
	Privacy   *PrivacyManager
	Staking   *StakingManager
	Governance *GovernanceManager
}

// NewClient creates a new Chert client with the given configuration
func NewClient(config *ClientConfig) (*ChertClient, error) {
	if config == nil {
		config = DefaultClientConfig()
	}

	if config.Endpoint == "" {
		config.Endpoint = DefaultEndpoint
	}

	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	if config.Network == "" {
		config.Network = NetworkMainnet
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	rpcClient := NewRPCClient(config.Endpoint, config.Timeout)

	client := &ChertClient{
		config:     config,
		httpClient: httpClient,
		rpcClient:  rpcClient,
	}

	// Initialize managers
	client.Wallet = NewWalletManager(client)
	client.Privacy = NewPrivacyManager(client)
	client.Staking = NewStakingManager(client)
	client.Governance = NewGovernanceManager(client)

	return client, nil
}

// GetNetworkStatus retrieves the current network status
func (c *ChertClient) GetNetworkStatus(ctx context.Context) (*NetworkStatus, error) {
	var result NetworkStatus
	err := c.rpcClient.Call(ctx, "getNetworkStatus", nil, &result)
	return &result, err
}

// GetLatestBlock retrieves the latest block information
func (c *ChertClient) GetLatestBlock(ctx context.Context) (*Block, error) {
	var result Block
	err := c.rpcClient.Call(ctx, "getLatestBlock", nil, &result)
	return &result, err
}

// GetBlock retrieves block information by height
func (c *ChertClient) GetBlock(ctx context.Context, height uint64) (*Block, error) {
	var result Block
	err := c.rpcClient.Call(ctx, "getBlock", []interface{}{height}, &result)
	return &result, err
}

// GetTransaction retrieves transaction information by hash
func (c *ChertClient) GetTransaction(ctx context.Context, hash string) (*Transaction, error) {
	var result Transaction
	err := c.rpcClient.Call(ctx, "getTransaction", []interface{}{hash}, &result)
	return &result, err
}

// IsConnected checks if the client is connected to the network
func (c *ChertClient) IsConnected(ctx context.Context) bool {
	_, err := c.GetNetworkStatus(ctx)
	return err == nil
}

// GetConfig returns the client configuration
func (c *ChertClient) GetConfig() *ClientConfig {
	return c.config
}

// makeRequest performs an HTTP request with proper error handling
func (c *ChertClient) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	url := c.config.Endpoint + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add API key if provided
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	// Add custom headers
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	return c.httpClient.Do(req)
}

// handleResponse processes an HTTP response and unmarshals the result
func (c *ChertClient) handleResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return &apiErr
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !apiResp.Success {
		if apiResp.Error != nil {
			return apiResp.Error
		}
		return fmt.Errorf("API request failed")
	}

	if result != nil {
		resultBytes, err := json.Marshal(apiResp.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal result data: %w", err)
		}
		return json.Unmarshal(resultBytes, result)
	}

	return nil
}

// APIResponse represents a standard API response
type APIResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %s: %s", e.Code, e.Message)
}

// GenerateAddress generates a deterministic address from a public key
func GenerateAddress(publicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("invalid public key hex: %w", err)
	}

	hash := sha256.Sum256(pubKeyBytes)
	address := "chert_" + hex.EncodeToString(hash[:20])
	return address, nil
}

// GenerateTxID generates a new transaction ID
func GenerateTxID() string {
	return uuid.New().String()
}
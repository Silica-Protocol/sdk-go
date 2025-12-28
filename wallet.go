package chert

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// WalletManager handles wallet operations and account management
type WalletManager struct {
	client *ChertClient
}

// NewWalletManager creates a new wallet manager
func NewWalletManager(client *ChertClient) *WalletManager {
	return &WalletManager{client: client}
}

// CreateAccount creates a new account with a randomly generated keypair
func (wm *WalletManager) CreateAccount() (*Account, error) {
	privateKey, publicKey, err := wm.generateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	address, err := GenerateAddress(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &Account{
		Address:    address,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// ImportAccount imports an account from a private key
func (wm *WalletManager) ImportAccount(privateKey string) (*Account, error) {
	publicKey, err := wm.derivePublicKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive public key: %w", err)
	}

	address, err := GenerateAddress(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &Account{
		Address:    address,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// CreateWatchOnlyAccount creates a watch-only account from a public key
func (wm *WalletManager) CreateWatchOnlyAccount(publicKey string) (*Account, error) {
	address, err := GenerateAddress(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &Account{
		Address:   address,
		PublicKey: publicKey,
	}, nil
}

// GetBalance retrieves the balance for an account
func (wm *WalletManager) GetBalance(ctx context.Context, address string) (*Balance, error) {
	var result Balance
	err := wm.client.rpcClient.Call(ctx, "getBalance", []interface{}{address}, &result)
	return &result, err
}

// SendTransaction sends a transaction to the network
func (wm *WalletManager) SendTransaction(ctx context.Context, request *TransactionRequest, account *Account) (string, error) {
	if account.PrivateKey == "" {
		return "", fmt.Errorf("account does not have a private key")
	}

	// Sign the transaction
	signature, err := wm.signTransaction(request, account.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	tx := map[string]interface{}{
		"sender":    account.Address,
		"recipient": request.To,
		"amount":    request.Amount,
		"fee":       request.Fee,
		"nonce":     request.Nonce,
		"signature": signature,
	}

	if request.Memo != "" {
		tx["memo"] = request.Memo
	}

	var result map[string]interface{}
	err = wm.client.rpcClient.Call(ctx, "sendTransaction", []interface{}{tx}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid transaction response")
}

// EstimateFee estimates the fee for a transaction
func (wm *WalletManager) EstimateFee(ctx context.Context, request *TransactionRequest) (*Fee, error) {
	var result Fee
	err := wm.client.rpcClient.Call(ctx, "estimateFee", []interface{}{request}, &result)
	return &result, err
}

// WaitForTransaction waits for a transaction to be confirmed
func (wm *WalletManager) WaitForTransaction(ctx context.Context, txHash string, timeoutMs uint64) (*Transaction, error) {
	// Default timeout of 60 seconds
	if timeoutMs == 0 {
		timeoutMs = 60000
	}

	interval := uint64(2000) // 2 seconds
	startTime := uint64(0)   // Would use time.Now() in real implementation

	for startTime < timeoutMs {
		tx, err := wm.client.GetTransaction(ctx, txHash)
		if err != nil {
			// Continue polling if transaction not found yet
			continue
		}

		if tx.Status == string(TxStatusConfirmed) {
			return tx, nil
		}

		if tx.Status == string(TxStatusFailed) || tx.Status == string(TxStatusRejected) {
			return nil, fmt.Errorf("transaction %s", tx.Status)
		}

		// Wait before next poll
		// In real implementation, would use time.Sleep
		startTime += interval
	}

	return nil, fmt.Errorf("transaction confirmation timeout")
}

// generateKeyPair generates a new Ed25519 keypair
func (wm *WalletManager) generateKeyPair() (string, string, error) {
	// Generate 32 bytes of random data for private key
	privateKeyBytes := make([]byte, 32)
	if _, err := rand.Read(privateKeyBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// In a real implementation, this would derive the public key from the private key
	// For now, we'll just use the private key bytes as public key for simplicity
	publicKeyBytes := make([]byte, 32)
	copy(publicKeyBytes, privateKeyBytes)

	privateKey := hex.EncodeToString(privateKeyBytes)
	publicKey := hex.EncodeToString(publicKeyBytes)

	return privateKey, publicKey, nil
}

// derivePublicKey derives the public key from a private key
func (wm *WalletManager) derivePublicKey(privateKey string) (string, error) {
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	if len(privateKeyBytes) != 32 {
		return "", fmt.Errorf("invalid private key length")
	}

	// In a real implementation, this would properly derive the public key
	// For now, we'll just return the private key bytes as public key
	publicKey := hex.EncodeToString(privateKeyBytes)
	return publicKey, nil
}

// signTransaction signs a transaction with the private key
func (wm *WalletManager) signTransaction(request *TransactionRequest, privateKey string) (string, error) {
	// In a real implementation, this would create a proper transaction hash
	// and sign it with the private key using Ed25519
	// For now, we'll return a mock signature
	return "mock_signature_" + GenerateTxID(), nil
}
# Chert SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/silica-network/chert/sdk/go.svg)](https://pkg.go.dev/github.com/silica-network/chert/sdk/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/silica-network/chert/sdk/go)](https://goreportcard.com/report/github.com/silica-network/chert/sdk/go)

Official Go SDK for the Chert/Silica blockchain network.

## Features

- **Wallet Management**: Create and manage accounts, send transactions
- **Privacy Features**: Stealth addresses and private transactions
- **Staking**: Delegate tokens to validators and manage stakes
- **Governance**: Participate in network governance and voting
- **Network Operations**: Query blockchain state and network information
- **Concurrent**: Built with Go's concurrency features
- **Type Safe**: Full type safety with Go's type system

## Installation

```bash
go get github.com/silica-network/chert/sdk/go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    chert "github.com/silica-network/chert/sdk/go"
)

func main() {
    // Create a client
    client, err := chert.NewClient(&chert.ClientConfig{
        Endpoint: "https://api.chert.com",
        Network:  chert.NetworkMainnet,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a new account
    account, err := client.Wallet.CreateAccount()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created account: %s\n", account.Address)

    // Get network status
    ctx := context.Background()
    status, err := client.GetNetworkStatus(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Current block height: %d\n", status.BlockHeight)
}
```

## Usage Examples

### Wallet Operations

```go
// Create a new account
account, err := client.Wallet.CreateAccount()
if err != nil {
    log.Fatal(err)
}

// Import account from private key
importedAccount, err := client.Wallet.ImportAccount("your_private_key_here")
if err != nil {
    log.Fatal(err)
}

// Get account balance
balance, err := client.Wallet.GetBalance(ctx, account.Address)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Balance: %s available, %s total\n", balance.Available, balance.Total)

// Send a transaction
txRequest := &chert.TransactionRequest{
    To:     "recipient_address",
    Amount: "100.0",
    Fee:    "0.1",
    Memo:   "Hello Chert!",
}

txHash, err := client.Wallet.SendTransaction(ctx, txRequest, account)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Transaction sent: %s\n", txHash)
```

### Privacy Features

```go
// Generate stealth keys
stealthKeys, err := client.Privacy.GenerateStealthKeys()
if err != nil {
    log.Fatal(err)
}

// Create stealth account
stealthAccount, err := client.Privacy.CreateStealthAccount(
    stealthKeys.ViewKeypair.Secret,
    stealthKeys.SpendKeypair.Public,
    stealthKeys,
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Stealth address: %s\n", stealthAccount.Address)

// Send private transaction
privateTxRequest := &chert.PrivateTransactionRequest{
    SenderKeys:   *stealthKeys,
    RecipientViewKey: "recipient_view_key",
    Amount:       "50.0",
    Fee:          "0.05",
    PrivacyLevel: chert.PrivacyLevelStealth,
    Memo:         "Private transaction",
}

privateTxID, err := client.Privacy.SendPrivateTransaction(ctx, privateTxRequest, "recipient_view_key", "recipient_spend_key")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Private transaction sent: %s\n", privateTxID)
```

### Staking Operations

```go
// Get validators
validators, err := client.Staking.GetValidators(ctx)
if err != nil {
    log.Fatal(err)
}

for _, validator := range validators {
    fmt.Printf("Validator: %s (%s)\n", validator.Name, validator.Status)
}

// Delegate tokens
delegationTx, err := client.Staking.Delegate(ctx, account.Address, validatorAddress, "1000.0", "0.1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Delegation transaction: %s\n", delegationTx)

// Get staking rewards
rewards, err := client.Staking.GetStakingRewards(ctx, account.Address)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Available rewards: %s\n", rewards.Available)
```

### Governance Operations

```go
// Get proposals
proposals, err := client.Governance.GetProposals(ctx, 10)
if err != nil {
    log.Fatal(err)
}

for _, proposal := range proposals {
    fmt.Printf("Proposal: %s (%s)\n", proposal.Title, proposal.Status)
}

// Create a proposal
proposalID, err := client.Governance.CreateProposal(ctx,
    "Network Upgrade Proposal",
    "Proposal to upgrade the network to version 2.0",
    account.Address,
    "1.0",
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Proposal created: %s\n", proposalID)

// Vote on proposal
voteTx, err := client.Governance.Vote(ctx, proposalID, account.Address, chert.VoteOptionYes, "0.1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Vote cast: %s\n", voteTx)
```

## Configuration

```go
config := &chert.ClientConfig{
    Endpoint: "https://api.chert.com",
    Network:  chert.NetworkTestnet,
    Timeout:  45 * time.Second, // 45 second timeout
    APIKey:   "your_api_key_here",
    Headers: map[string]string{
        "X-Custom-Header": "value",
    },
}

client, err := chert.NewClient(config)
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

The SDK provides comprehensive error handling:

```go
balance, err := client.Wallet.GetBalance(ctx, address)
if err != nil {
    switch e := err.(type) {
    case *chert.APIError:
        fmt.Printf("API error %s: %s\n", e.Code, e.Message)
    default:
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

## Network Support

- **Mainnet**: Production network
- **Testnet**: Testing network for development
- **Devnet**: Local development network

## API Reference

For complete API documentation, see [pkg.go.dev/github.com/silica-network/chert/sdk/go](https://pkg.go.dev/github.com/silica-network/chert/sdk/go).

## Testing

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please see our [contributing guidelines](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

This SDK handles cryptographic operations and private keys. Always:

- Use strong, randomly generated private keys
- Never log or expose private keys
- Use HTTPS endpoints in production
- Keep dependencies updated
- Audit your code for security vulnerabilities
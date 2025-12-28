// Basic usage example for the Chert SDK
//
// This example demonstrates:
// - Creating a client
// - Creating and managing accounts
// - Sending transactions
// - Querying blockchain state
// - Basic error handling

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	chert "github.com/silica-network/chert/sdk/go"
)

func main() {
	fmt.Println("🚀 Chert SDK Go Basic Usage Example")
	fmt.Println("===================================")

	// Create a client with default configuration
	fmt.Println("\n📡 Creating Chert client...")
	client, err := chert.NewClient(chert.DefaultClientConfig())
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fmt.Println("✅ Client created successfully")

	// Test network connectivity
	fmt.Println("\n🌐 Testing network connectivity...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client.IsConnected(ctx) {
		fmt.Println("✅ Connected to Chert network")

		// Get network status
		fmt.Println("\n📊 Getting network status...")
		status, err := client.GetNetworkStatus(ctx)
		if err != nil {
			fmt.Printf("❌ Failed to get network status: %v\n", err)
		} else {
			fmt.Printf("✅ Network status:\n")
			fmt.Printf("   Block height: %d\n", status.BlockHeight)
			fmt.Printf("   Network ID: %s\n", status.NetworkID)
			fmt.Printf("   Consensus version: %s\n", status.ConsensusVersion)
			fmt.Printf("   Connected peers: %d\n", status.PeerCount)
			fmt.Printf("   Syncing: %t\n", status.Syncing)
		}

		// Get latest block
		fmt.Println("\n📦 Getting latest block...")
		block, err := client.GetLatestBlock(ctx)
		if err != nil {
			fmt.Printf("❌ Failed to get latest block: %v\n", err)
		} else {
			fmt.Printf("✅ Latest block:\n")
			fmt.Printf("   Height: %d\n", block.Height)
			fmt.Printf("   Hash: %s\n", block.Hash)
			fmt.Printf("   Transactions: %d\n", block.TransactionCount)
			fmt.Printf("   Proposer: %s\n", block.Proposer)
		}

		// Get validators (staking)
		fmt.Println("\n🏛️  Getting validators...")
		validators, err := client.Staking.GetValidators(ctx)
		if err != nil {
			fmt.Printf("❌ Failed to get validators: %v\n", err)
		} else {
			fmt.Printf("✅ Found %d validators\n", len(validators))
			for i, validator := range validators {
				if i >= 3 { // Show only first 3
					break
				}
				fmt.Printf("   Validator %d: %s (%s)\n", i+1, validator.Name, validator.Status)
			}
			if len(validators) > 3 {
				fmt.Printf("   ... and %d more\n", len(validators)-3)
			}
		}

		// Get governance proposals
		fmt.Println("\n🗳️  Getting governance proposals...")
		proposals, err := client.Governance.GetProposals(ctx, 10)
		if err != nil {
			fmt.Printf("❌ Failed to get proposals: %v\n", err)
		} else {
			fmt.Printf("✅ Found %d proposals\n", len(proposals))
			for i, proposal := range proposals {
				if i >= 3 { // Show only first 3
					break
				}
				fmt.Printf("   Proposal %d: %s (%s)\n", i+1, proposal.Title, proposal.Status)
			}
		}

	} else {
		fmt.Println("❌ Not connected to network - running offline examples")
		offlineExamples(client)
		return
	}

	// Create a new account
	fmt.Println("\n👛 Creating new account...")
	account, err := client.Wallet.CreateAccount()
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}

	fmt.Printf("✅ Account created:\n")
	fmt.Printf("   Address: %s\n", account.Address)
	fmt.Printf("   Public key: %s\n", account.PublicKey)

	// Import account from private key (example)
	fmt.Println("\n🔑 Importing account from private key...")
	examplePrivateKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	importedAccount, err := client.Wallet.ImportAccount(examplePrivateKey)
	if err != nil {
		fmt.Printf("❌ Failed to import account: %v\n", err)
	} else {
		fmt.Printf("✅ Account imported:\n")
		fmt.Printf("   Address: %s\n", importedAccount.Address)
	}

	// Get account balance (will likely fail without real account)
	fmt.Println("\n💰 Checking account balance...")
	balance, err := client.Wallet.GetBalance(ctx, account.Address)
	if err != nil {
		fmt.Printf("❌ Failed to get balance (expected for new account): %v\n", err)
	} else {
		fmt.Printf("✅ Account balance:\n")
		fmt.Printf("   Available: %s\n", balance.Available)
		fmt.Printf("   Pending: %s\n", balance.Pending)
		fmt.Printf("   Total: %s\n", balance.Total)
	}

	// Estimate transaction fee
	fmt.Println("\n💸 Estimating transaction fee...")
	txRequest := &chert.TransactionRequest{
		To:     "chert_1example_address_1234567890",
		Amount: "100.0",
		Fee:    "0.1",
		Memo:   "Chert SDK Go Example Transaction",
	}

	fee, err := client.Wallet.EstimateFee(ctx, txRequest)
	if err != nil {
		fmt.Printf("❌ Failed to estimate fee: %v\n", err)
	} else {
		fmt.Printf("✅ Fee estimate:\n")
		fmt.Printf("   Amount: %s\n", fee.Amount)
		if fee.GasLimit > 0 {
			fmt.Printf("   Gas limit: %d\n", fee.GasLimit)
		}
	}

	fmt.Println("\n🎉 Example completed successfully!")
	fmt.Println("💡 Tip: Use testnet for development and mainnet for production")
}

// Examples that work without network connectivity
func offlineExamples(client *chert.ChertClient) {
	fmt.Println("\n🏠 Running offline examples...")

	// Account creation (local operation)
	fmt.Println("\n👛 Creating accounts offline...")
	account1, err := client.Wallet.CreateAccount()
	if err != nil {
		log.Fatal("Failed to create account1:", err)
	}

	account2, err := client.Wallet.CreateAccount()
	if err != nil {
		log.Fatal("Failed to create account2:", err)
	}

	fmt.Printf("✅ Account 1: %s\n", account1.Address)
	fmt.Printf("✅ Account 2: %s\n", account2.Address)

	// Create transaction request (doesn't send)
	fmt.Println("\n📝 Creating transaction request...")
	txRequest := &chert.TransactionRequest{
		To:     account2.Address,
		Amount: "50.0",
		Fee:    "0.05",
		Memo:   "Offline example transaction",
	}

	fmt.Printf("✅ Transaction request created:\n")
	fmt.Printf("   From: %s\n", account1.Address)
	fmt.Printf("   To: %s\n", txRequest.To)
	fmt.Printf("   Amount: %s\n", txRequest.Amount)
	fmt.Printf("   Fee: %s\n", txRequest.Fee)
	fmt.Printf("   Memo: %s\n", txRequest.Memo)

	// Privacy features (local)
	fmt.Println("\n🔒 Generating stealth keys...")
	stealthKeys, err := client.Privacy.GenerateStealthKeys()
	if err != nil {
		log.Fatal("Failed to generate stealth keys:", err)
	}

	fmt.Printf("✅ Stealth keys generated:\n")
	fmt.Printf("   View public: %s\n", stealthKeys.ViewKeypair.Public)
	fmt.Printf("   Spend public: %s\n", stealthKeys.SpendKeypair.Public)

	// Create stealth account
	stealthAccount, err := client.Privacy.CreateStealthAccount(
		stealthKeys.ViewKeypair.Secret,
		stealthKeys.SpendKeypair.Public,
		stealthKeys,
	)
	if err != nil {
		log.Fatal("Failed to create stealth account:", err)
	}

	fmt.Printf("✅ Stealth account created:\n")
	fmt.Printf("   Address: %s\n", stealthAccount.Address)

	fmt.Println("\n🎉 Offline examples completed!")
	fmt.Println("💡 To use network features, ensure you have a valid Chert API endpoint")
}
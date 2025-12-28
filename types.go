package chert

import (
	"fmt"
	"time"
)

// Account represents a blockchain account
type Account struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key,omitempty"`
}

// Balance represents an account balance
type Balance struct {
	Available string `json:"available"`
	Pending   string `json:"pending"`
	Total     string `json:"total"`
}

// TransactionRequest represents a transaction request
type TransactionRequest struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
	Fee    string `json:"fee"`
	Memo   string `json:"memo,omitempty"`
	Nonce  uint64 `json:"nonce,omitempty"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Amount      string    `json:"amount"`
	Fee         string    `json:"fee"`
	Memo        string    `json:"memo,omitempty"`
	BlockHeight uint64    `json:"block_height,omitempty"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Nonce       uint64    `json:"nonce"`
}

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TxStatusPending   TransactionStatus = "pending"
	TxStatusConfirmed TransactionStatus = "confirmed"
	TxStatusFailed    TransactionStatus = "failed"
	TxStatusRejected  TransactionStatus = "rejected"
)

// Privacy types
type StealthKeys struct {
	ViewKeypair  KeyPair `json:"view_keypair"`
	SpendKeypair KeyPair `json:"spend_keypair"`
}

type KeyPair struct {
	Public  string `json:"public"`
	Secret  string `json:"secret"`
}

type StealthAccount struct {
	Address        string       `json:"address"`
	ViewKey        string       `json:"view_key"`
	SpendPublicKey string       `json:"spend_public_key"`
	Keys           *StealthKeys `json:"keys,omitempty"`
}

type PrivacyLevel string

const (
	PrivacyLevelStealth  PrivacyLevel = "stealth"
	PrivacyLevelEncrypted PrivacyLevel = "encrypted"
)

type PrivateTransactionRequest struct {
	SenderKeys    StealthKeys  `json:"sender_keys"`
	RecipientViewKey string    `json:"recipient_view_key"`
	Amount        string       `json:"amount"`
	Fee           string       `json:"fee"`
	Memo          string       `json:"memo,omitempty"`
	PrivacyLevel  PrivacyLevel `json:"privacy_level"`
	Nonce         uint64       `json:"nonce"`
}

type PrivateTransaction struct {
	TxID      string    `json:"tx_id"`
	Amount    string    `json:"amount"`
	Memo      string    `json:"memo,omitempty"`
	Sender    string    `json:"sender,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Fee       string    `json:"fee"`
}

// Staking types
type Validator struct {
	Address         string `json:"address"`
	Name            string `json:"name"`
	VotingPower     string `json:"voting_power"`
	Commission      string `json:"commission"`
	Status          string `json:"status"`
	TotalDelegated  string `json:"total_delegated"`
	DelegatorCount  uint64 `json:"delegator_count"`
	PublicKey       string `json:"public_key,omitempty"`
	StakeAmount     uint64 `json:"stake_amount,omitempty"`
	CommissionRate  uint32 `json:"commission_rate,omitempty"`
	IsActive        bool   `json:"is_active,omitempty"`
	ReputationScore float64 `json:"reputation_score,omitempty"`
	LastActivity    time.Time `json:"last_activity,omitempty"`
}

type ValidatorStatus string

const (
	ValidatorStatusActive   ValidatorStatus = "active"
	ValidatorStatusInactive ValidatorStatus = "inactive"
	ValidatorStatusJailed   ValidatorStatus = "jailed"
)

type DelegationRequest struct {
	ValidatorAddress string `json:"validator_address"`
	Amount           string `json:"amount"`
	Fee              string `json:"fee"`
}

type Delegation struct {
	ValidatorAddress string    `json:"validator_address"`
	Amount           string    `json:"amount"`
	Rewards          string    `json:"rewards"`
	Timestamp        time.Time `json:"timestamp"`
}

type StakingRewards struct {
	Total    string     `json:"total"`
	Available string    `json:"available"`
	Pending   string     `json:"pending"`
	LastClaim *time.Time `json:"last_claim,omitempty"`
}

// Governance types
type Proposal struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Proposer        string    `json:"proposer"`
	Status          string    `json:"status"`
	VotingStartTime time.Time `json:"voting_start_time"`
	VotingEndTime   time.Time `json:"voting_end_time"`
	Tally           VoteTally `json:"tally"`
}

type ProposalStatus string

const (
	ProposalStatusVoting    ProposalStatus = "voting"
	ProposalStatusPassed    ProposalStatus = "passed"
	ProposalStatusRejected  ProposalStatus = "rejected"
	ProposalStatusExecuted  ProposalStatus = "executed"
	ProposalStatusFailed    ProposalStatus = "failed"
)

type VoteTally struct {
	Yes        string `json:"yes"`
	No         string `json:"no"`
	Abstain    string `json:"abstain"`
	NoWithVeto string `json:"no_with_veto"`
}

type VoteOption string

const (
	VoteOptionYes        VoteOption = "yes"
	VoteOptionNo         VoteOption = "no"
	VoteOptionAbstain    VoteOption = "abstain"
	VoteOptionNoWithVeto VoteOption = "no_with_veto"
)

type VoteRequest struct {
	ProposalID string    `json:"proposal_id"`
	Option     VoteOption `json:"option"`
	Fee        string     `json:"fee"`
}

// Network types
type NetworkStatus struct {
	BlockHeight      uint64    `json:"block_height"`
	NetworkID        string    `json:"network_id"`
	ConsensusVersion string    `json:"consensus_version"`
	PeerCount        uint64    `json:"peer_count"`
	Syncing          bool      `json:"syncing"`
	LatestBlockTime  time.Time `json:"latest_block_time"`
}

type Block struct {
	Height           uint64        `json:"height"`
	Hash             string        `json:"hash"`
	PreviousHash     string        `json:"previous_hash"`
	Timestamp        time.Time     `json:"timestamp"`
	TransactionCount uint64        `json:"transaction_count"`
	Proposer         string        `json:"proposer"`
	Transactions     []Transaction `json:"transactions,omitempty"`
}

// Fee estimation
type Fee struct {
	Amount    string `json:"amount"`
	GasLimit  uint64 `json:"gas_limit,omitempty"`
	GasPrice  string `json:"gas_price,omitempty"`
}

// JSON-RPC types
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  interface{}   `json:"params,omitempty"`
	ID      interface{}   `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}
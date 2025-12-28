package chert

import (
	"context"
	"fmt"
)

// GovernanceManager handles governance operations and proposals
type GovernanceManager struct {
	client *ChertClient
}

// NewGovernanceManager creates a new governance manager
func NewGovernanceManager(client *ChertClient) *GovernanceManager {
	return &GovernanceManager{client: client}
}

// GetProposals retrieves the list of governance proposals
func (gm *GovernanceManager) GetProposals(ctx context.Context, limit int) ([]*Proposal, error) {
	params := make(map[string]interface{})
	if limit > 0 {
		params["limit"] = limit
	}

	var result struct {
		Proposals []*Proposal `json:"proposals"`
	}
	err := gm.client.rpcClient.Call(ctx, "governance_getProposals", []interface{}{params}, &result)
	return result.Proposals, err
}

// GetProposal retrieves a specific proposal by ID
func (gm *GovernanceManager) GetProposal(ctx context.Context, proposalID string) (*Proposal, error) {
	var result Proposal
	err := gm.client.rpcClient.Call(ctx, "governance_getProposal", []interface{}{proposalID}, &result)
	return &result, err
}

// CreateProposal creates a new governance proposal
func (gm *GovernanceManager) CreateProposal(ctx context.Context, title, description, proposerAddress, fee string) (string, error) {
	params := map[string]interface{}{
		"title":       title,
		"description": description,
		"proposer":    proposerAddress,
		"fee":         fee,
	}

	var result map[string]interface{}
	err := gm.client.rpcClient.Call(ctx, "governance_createProposal", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if proposalID, ok := result["proposal_id"].(string); ok {
		return proposalID, nil
	}

	return "", fmt.Errorf("invalid proposal creation response")
}

// Vote casts a vote on a governance proposal
func (gm *GovernanceManager) Vote(ctx context.Context, proposalID, voterAddress string, option VoteOption, fee string) (string, error) {
	params := map[string]interface{}{
		"proposal_id": proposalID,
		"voter":       voterAddress,
		"option":      option,
		"fee":         fee,
	}

	var result map[string]interface{}
	err := gm.client.rpcClient.Call(ctx, "governance_vote", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid vote response")
}

// GetProposalVotes retrieves votes for a specific proposal
func (gm *GovernanceManager) GetProposalVotes(ctx context.Context, proposalID string) (*VoteTally, error) {
	var result VoteTally
	err := gm.client.rpcClient.Call(ctx, "governance_getProposalVotes", []interface{}{proposalID}, &result)
	return &result, err
}

// GetVoterVotes retrieves votes cast by a specific voter
func (gm *GovernanceManager) GetVoterVotes(ctx context.Context, voterAddress string) (map[string]VoteOption, error) {
	var result map[string]VoteOption
	err := gm.client.rpcClient.Call(ctx, "governance_getVoterVotes", []interface{}{voterAddress}, &result)
	return result, err
}

// ExecuteProposal executes a passed proposal (admin function)
func (gm *GovernanceManager) ExecuteProposal(ctx context.Context, proposalID, executorAddress, fee string) (string, error) {
	params := map[string]interface{}{
		"proposal_id": proposalID,
		"executor":    executorAddress,
		"fee":         fee,
	}

	var result map[string]interface{}
	err := gm.client.rpcClient.Call(ctx, "governance_executeProposal", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid proposal execution response")
}

// CancelProposal cancels a proposal (only by proposer)
func (gm *GovernanceManager) CancelProposal(ctx context.Context, proposalID, proposerAddress, fee string) (string, error) {
	params := map[string]interface{}{
		"proposal_id": proposalID,
		"proposer":    proposerAddress,
		"fee":         fee,
	}

	var result map[string]interface{}
	err := gm.client.rpcClient.Call(ctx, "governance_cancelProposal", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid proposal cancellation response")
}

// GetProposalStatus retrieves the current status of a proposal
func (gm *GovernanceManager) GetProposalStatus(ctx context.Context, proposalID string) (ProposalStatus, error) {
	var result struct {
		Status ProposalStatus `json:"status"`
	}
	err := gm.client.rpcClient.Call(ctx, "governance_getProposalStatus", []interface{}{proposalID}, &result)
	return result.Status, err
}

// GetVotingPower retrieves the voting power of an address
func (gm *GovernanceManager) GetVotingPower(ctx context.Context, address string) (string, error) {
	var result struct {
		VotingPower string `json:"voting_power"`
	}
	err := gm.client.rpcClient.Call(ctx, "governance_getVotingPower", []interface{}{address}, &result)
	return result.VotingPower, err
}

// GetGovernanceStats retrieves governance statistics
func (gm *GovernanceManager) GetGovernanceStats(ctx context.Context) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := gm.client.rpcClient.Call(ctx, "governance_getStats", nil, &result)
	return result, err
}
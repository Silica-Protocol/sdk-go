package chert

import (
	"context"
	"fmt"
)

// StakingManager handles staking and delegation operations
type StakingManager struct {
	client *ChertClient
}

// NewStakingManager creates a new staking manager
func NewStakingManager(client *ChertClient) *StakingManager {
	return &StakingManager{client: client}
}

// GetValidators retrieves the list of validators
func (sm *StakingManager) GetValidators(ctx context.Context) ([]*Validator, error) {
	var result struct {
		Validators []*Validator `json:"validators"`
	}
	err := sm.client.rpcClient.Call(ctx, "getValidators", nil, &result)
	return result.Validators, err
}

// GetValidator retrieves a specific validator by address
func (sm *StakingManager) GetValidator(ctx context.Context, address string) (*Validator, error) {
	var result Validator
	err := sm.client.rpcClient.Call(ctx, "getValidator", []interface{}{address}, &result)
	return &result, err
}

// Delegate delegates tokens to a validator
func (sm *StakingManager) Delegate(ctx context.Context, delegatorAddress, validatorAddress, amount, fee string) (string, error) {
	params := map[string]interface{}{
		"delegator": delegatorAddress,
		"validator": validatorAddress,
		"amount":    amount,
		"fee":       fee,
	}

	var result map[string]interface{}
	err := sm.client.rpcClient.Call(ctx, "staking_delegate", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid delegation response")
}

// Undelegate removes delegation from a validator
func (sm *StakingManager) Undelegate(ctx context.Context, delegatorAddress, validatorAddress, amount, fee string) (string, error) {
	params := map[string]interface{}{
		"delegator": delegatorAddress,
		"validator": validatorAddress,
		"amount":    amount,
		"fee":       fee,
	}

	var result map[string]interface{}
	err := sm.client.rpcClient.Call(ctx, "staking_undelegate", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid undelegation response")
}

// GetDelegations retrieves delegations for an account
func (sm *StakingManager) GetDelegations(ctx context.Context, delegatorAddress string) ([]*Delegation, error) {
	var result struct {
		Delegations []*Delegation `json:"delegations"`
	}
	err := sm.client.rpcClient.Call(ctx, "getDelegations", []interface{}{delegatorAddress}, &result)
	return result.Delegations, err
}

// GetStakingRewards retrieves staking rewards for an account
func (sm *StakingManager) GetStakingRewards(ctx context.Context, delegatorAddress string) (*StakingRewards, error) {
	var result StakingRewards
	err := sm.client.rpcClient.Call(ctx, "getStakingRewards", []interface{}{delegatorAddress}, &result)
	return &result, err
}

// ClaimRewards claims staking rewards
func (sm *StakingManager) ClaimRewards(ctx context.Context, delegatorAddress string, validatorAddress, fee string) (string, error) {
	params := map[string]interface{}{
		"delegator": delegatorAddress,
		"validator": validatorAddress,
		"fee":       fee,
	}

	var result map[string]interface{}
	err := sm.client.rpcClient.Call(ctx, "staking_claimRewards", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid claim rewards response")
}

// RegisterValidator registers a new validator
func (sm *StakingManager) RegisterValidator(ctx context.Context, validator *Validator, ownerAddress, fee string) (string, error) {
	params := map[string]interface{}{
		"validator":     validator,
		"owner_address": ownerAddress,
		"fee":           fee,
	}

	var result map[string]interface{}
	err := sm.client.rpcClient.Call(ctx, "staking_registerValidator", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid validator registration response")
}

// UpdateCommission updates a validator's commission rate
func (sm *StakingManager) UpdateCommission(ctx context.Context, validatorAddress, ownerAddress string, newRate uint32, fee string) (string, error) {
	params := map[string]interface{}{
		"validator_address": validatorAddress,
		"owner_address":     ownerAddress,
		"new_rate":          newRate,
		"fee":               fee,
	}

	var result map[string]interface{}
	err := sm.client.rpcClient.Call(ctx, "staking_updateCommission", []interface{}{params}, &result)
	if err != nil {
		return "", err
	}

	if txHash, ok := result["tx_hash"].(string); ok {
		return txHash, nil
	}

	return "", fmt.Errorf("invalid commission update response")
}
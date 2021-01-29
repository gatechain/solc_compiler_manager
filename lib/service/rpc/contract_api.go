package rpc

import (
	"context"
	"github.com/gatechain/smart_contract_verifier/lib"
	"github.com/gatechain/smart_contract_verifier/lib/compiler/solidity"
	"log"
)

// ContractAPI is the contract_ prefixed set of APIs in the JSON-RPC spec.
type ContractAPI struct {
	cliCtx      context.Context
	logger      log.Logger
}

// NewContractAPI creates an instance of the Web3 API.
func NewContractAPI(ctx context.Context) *ContractAPI {
	return &ContractAPI{cliCtx: ctx}
}

// Ping server heart beat.
func (api *ContractAPI) Ping() string {
	return "Pong"
}

func (api *ContractAPI) Verify(params lib.CompileInput) (map[string]interface{}, error) {
	res, err := solidity.RemoteVerify(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}


package rpc

import (
	"context"
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

// Ping checks server is alive.
func (api *ContractAPI) Ping() string {
	return "pong"
}

func (api *ContractAPI) Verify() string {
	return "pong"
}

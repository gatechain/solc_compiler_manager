package rpc

import (
	"context"
	"log"
	"sync"
)

// ContractAPI is the contract_ prefixed set of APIs in the JSON-RPC spec.
type ContractAPI struct {
	cliCtx      context.Context
	logger      log.Logger
	keybaseLock sync.Mutex
}

// NewContractAPI creates an instance of the Web3 API.
func NewContractAPI() *ContractAPI {
	return &ContractAPI{}
}



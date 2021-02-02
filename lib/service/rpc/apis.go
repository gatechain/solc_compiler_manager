package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gatechain/solc_compiler_manager/lib/service/rest"
)

// RPC namespaces and API version
const (
	ContractNamespace     = "contract"
	apiVersion = "1.0"
)

// API describes the set of methods offered over the RPC interface
type API struct {
	Namespace string      // namespace under which the rpc methods of Service are exposed
	Version   string      // api version for DApp's
	Service   interface{} // receiver instance which holds the methods
	Public    bool        // indication if the methods must be considered safe for public use
}

func GetApis(cliCtx context.Context) []API{
	return []API{
		{
			Namespace: ContractNamespace,
			Version:   apiVersion,
			Service:   NewContractAPI(cliCtx),
			Public:    true,
		},
	}
}

// RegisterRoutes creates a new server and registers the `/rpc` endpoint.
// Rpc calls are enabled based on their associated module.
func RegisterRoutes(rs *rest.Server) {
	server := rpc.NewServer()
	apis := GetApis(rs.CliCtx)

	for _, api := range apis {
		if err := server.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}

	// register handler
	rs.Mux.HandleFunc("/", server.ServeHTTP).Methods("POST", "GET")
}

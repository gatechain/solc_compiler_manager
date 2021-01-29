module github.com/gatechain/smart_contract_verifier

go 1.14

require (
	github.com/bgentry/speakeasy v0.1.0
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.9.21
	github.com/go-kit/kit v0.8.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-isatty v0.0.5-0.20180830101745-3fb116b82035
	github.com/spf13/cobra v1.1.1
	github.com/tendermint/go-amino v0.16.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202
)

replace github.com/gatechain/smart_contract_verifier v0.0.0-20210125081747-cd368d6e121a => ./

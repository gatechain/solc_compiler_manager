package lib

// Solc compiler version
// SolcBuild and SolcVersion parse solc build version response
/* raw data
{
	"path": "soljson-v0.1.1+commit.6ff4cd6.js",
	"version": "0.1.1",
	"build": "commit.6ff4cd6",
	"longVersion": "0.1.1+commit.6ff4cd6",
	"keccak256": "0xd8b8c64f4e9de41e6604e6ac30274eff5b80f831f8534f0ad85ec0aff466bb25",
	"urls": [
		"bzzr://8f3c028825a1b72645f46920b67dca9432a87fc37a8940a2b2ce1dd6ddc2e29b",
		"dweb:/ipfs/QmPPGxsMtQSEUt9hn9VAtrtBTa1u9S5KF1myw78epNNFkx"
	]
}
*/
type SolcBuild struct {
	Path 		string 		`json:"path"`
	Version		string 		`json:"version"`
	Build		string 		`json:"build"`
	LongVersion string 		`json:"longVersion"`
	Keccak256	string 		`json:"keccak256"`
	Urls		[]string 	`json:"urls"`
}

type SolcVersion struct{
	Builds 			[]SolcBuild     	`json:"builds"`
	Releases		map[string]string 	`json:"releases"`
	LatestRelease 	string           	`json:"latestRelease"`
}

// solidity compile response
/*
@Input
{
	name: "SimpleStorage",
    compiler_version: "v0.4.24+commit.e67f0147",
    code: \"""
    pragma solidity ^0.4.24;

    contract SimpleStorage {
		uint storedData;

    	function set(uint x) public {
    		storedData = x;
    	}

		function get() public constant returns (uint) {
    		return storedData;
    	}
	}
    \""",
    optimize: false,
	evm_version: "byzantium"
}
@Output
{
	"abi": [
	  	{
			"constant" => false,
			"inputs" => [%{"name" => "x", "type" => "uint256"}],
			"name" => "set",
			"outputs" => [],
			"payable" => false,
			"stateMutability" => "nonpayable",
			"type" => "function"
	  	},
	  	{
			"constant" => true,
			"inputs" => [],
			"name" => "get",
			"outputs" => [%{"name" => "", "type" => "uint256"}],
			"payable" => false,
			"stateMutability" => "view",
			"type" => "function"
	  	},
	],
	"bytecode": "608060405234801561001057600080fd5b5060df8061001f6000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a72305820834bdab406d80509618957aa1a5ad1a4b77f4f1149078675940494ebe5b4147b0029",
	"name": "SimpleStorage"
}
*/
type CompileInput struct {
	Name            	string 	`json:"name"`
	CompilerVersion 	string 	`json:"compiler_version"`
	Code 				string 	`json:"code"`
	Optimize			bool	`json:"optimize"`
	OptimizationRuns	int		`json:"optimization_runs"`
	EvmVersion			string	`json:"evm_version"`
}

type ABI struct {
	ABI			[]string 		`json:"abi"`
	Bytecode	string   	`json:"bytecode"`
	Name		string   	`json:"name"`
}

type EventABI struct {
	Anonymous	bool 	`json:"anonymous"`
	Inputs		[]string	`json:"inputs"`
}

type EventABIIO struct {
	Indexed			bool 	`json:"indexed"`
	InternalType	string	`json:"internalType"`
	Name			string	`json:"name"`
}

type FunctionABI struct {
	Constant		bool      			`json:"constant"`
	Inputs			[]FunctionABIIO  	`json:"inputs"`
	Name			string    			`json:"name"`
	Outputs			[]FunctionABIIO 	`json:"outputs"`
	Payable			bool   				`json:"payable"`
	StateMutability	string 				`json:"stateMutability"`
	Type			string    			`json:"type"`
}

type FunctionABIIO struct {
	InternalType	string	`json:"internalType"`
	Name			string	`json:"name"`
	Type			string	`json:"type"`
}

type Bin string

type Hashes map[string]string

type Version string

// UnpackData handles local compile response, unmarshal json string to object
type UnpackData struct {
	Contracts map[string]map[string]interface{}	`json:"contracts"`
	Version										`json:"version"`
}
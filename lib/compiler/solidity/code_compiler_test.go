package solidity

import (
	"github.com/gatechain/solc_compiler_manager/lib"
	"testing"
)

func TestVerify(t *testing.T) {
	Input := lib.CompileInput{
		Name: "SimpleStorage",
		CompilerVersion: "v0.4.24+commit.e67f0147",
		Code: "pragma solidity ^0.4.24; contract SimpleStorage {uint storedData; function set(uint x) public { storedData = x;} function get() public constant returns (uint) { return storedData;} }",
		Optimize: false,
		OptimizationRuns: 200,
		EvmVersion: "byzantium",
	}
	res, err := RemoteVerify(Input)
	if err != nil {
		t.Errorf("Remote verify error: %s", err)
	}
	print(res)
}

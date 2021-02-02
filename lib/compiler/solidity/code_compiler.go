package solidity

import (
	"encoding/json"
	"fmt"
	"github.com/gatechain/solc_compiler_manager/lib"
	"github.com/gatechain/solc_compiler_manager/lib/compiler"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	//"math/rand"
	//"os"
	//"os/exec"
	"strings"
)

const allowedEvmVersion = "homestead,tangerineWhistle,spuriousDragon,byzantium,constantinople,petersburg,istanbul,default"
const NewContractName = "New.sol"
const FlagScope = "scope"
const FlagName = "name"
const FlagEvmVersion = "evm-version"
const FlagOptimize = "optimize"
const FlagOptimizeRuns = "optimize-runs"
const ScopeVerify = "bin,abi"
const ScopeABI = "abi"
const ScopeHashes = "hashes"

// Module responsible to compile the Solidity code of a given Smart Contract.
/* raw data
Compiles a code in the solidity command line.

Returns a `Map`.

## Examples

    ...> RemoteVerify({
    ...>      name: "SimpleStorage",
    ...>      compiler_version: "v0.4.24+commit.e67f0147",
    ...>      code: \"""
    ...>      pragma solidity ^0.4.24;
    ...>
    ...>      contract SimpleStorage {
    ...>          uint storedData;
    ...>
    ...>          function set(uint x) public {
    ...>              storedData = x;
    ...>          }
    ...>
    ...>          function get() public constant returns (uint) {
    ...>              return storedData;
    ...>          }
    ...>      }
    ...>      \""",
    ...>      optimize: false,
	...>	  optimize_runs: 0,
	...>	  evm_version: "byzantium"
    ...>  })
    return:
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
          	}
		],
        "bytecode" => "608060405234801561001057600080fd5b5060df8061001f6000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a72305820834bdab406d80509618957aa1a5ad1a4b77f4f1149078675940494ebe5b4147b0029",
        "name" => "SimpleStorage"
	}
*/
func RemoteVerify(params lib.CompileInput) (map[string]interface{}, error) {
	// check version and commit
	version, commit, err := lib.CheckLongVersionFormat(params.CompilerVersion)
	if err != nil {
		return nil, err
	}
	if !lib.CheckVersionCommit(version, commit) {
		return nil, fmt.Errorf("version with given commit not match, please check your compiler version")
	}

	// check evm version
	var evmVersion string
	if params.EvmVersion == "" {
		evmVersion = latestEvmVersion()
	} else {
		evmVersion = params.EvmVersion
	}
	checkedEvmVersion, _ := isEvmVersionAllowed(evmVersion)

	// solidity execute binary path
	executePath, err := compiler.EnsureExists(version)
	if err != nil {
		return nil, err
	}

	filePath, err := createSourceFile([]byte(params.Code))
	if err != nil {
		return nil, err
	}
	defer deleteSourceFile(filePath)

	res, err := LocalRun(
		executePath, filePath, params.Name, ScopeVerify,
		checkedEvmVersion, params.Optimize, params.OptimizationRuns,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func LocalRun(compilePath, filePath, name, scopes, evmVersion string, optimize bool, optimizeRuns int) (map[string]interface{}, error) {
	// format evm version
	if evmVersion == "default" {
		evmVersion = ""
	} else {
		evmVersion = fmt.Sprintf(" --evm-version %s", evmVersion)
	}
	// format option
	var optimizeFlag, runFlag string
	if optimize {
		optimizeFlag = " --optimize"
		runFlag = fmt.Sprintf(" --optimize-runs %d", optimizeRuns)
	} else {
		optimizeFlag = ""
		runFlag = ""
	}
	// execute compile
	cmd := fmt.Sprintf(
		"%s --combined-json %s --pretty-json %s",
		compilePath, scopes, filePath,
	) + evmVersion + optimizeFlag + runFlag

	command := exec.Command("bash", "-c", cmd)
	output, err := command.Output()
	if err != nil {
		return nil, err
	}
	// unpack output
	res, err := unpack(name, output)
	if err != nil {
		return nil, err
	}
	// print format
	for k, v := range res {
		fmt.Println(k, ": ")
		switch rv := v.(type) {
		case string:
			fmt.Println(v)
		case map[string]interface{}:
			for i, j := range rv {
				fmt.Println(j, "\t", i)
			}
		}
	}
	return res, nil
}

func unpack(name string, output []byte) (map[string]interface{}, error) {
	// unpack data
	tmp := lib.UnpackData{}
	err := json.Unmarshal(output, &tmp)
	if err != nil {
		return nil, fmt.Errorf("unpack data error")
	}
	// unpack output by given contract name
	Contract :=  make(map[string]interface{})
	for k, v := range tmp.Contracts {
		seps := strings.Split(k, ":")
		if len(seps) != 2 {
			return nil, fmt.Errorf("unpack data error")
		}
		if name == seps[1] {
			Contract = v
		}
	}
	if len(Contract) == 0 {
		return nil, fmt.Errorf("contract name not found")
	}

	return Contract, nil
}

func latestEvmVersion() string {
	version := allowedEvmVersions()
	return version[len(version) - 1]
}

func allowedEvmVersions() []string {
	return strings.Split(allowedEvmVersion, ",")
}

func isEvmVersionAllowed(evmVersion string) (string, bool) {
	if strings.Contains(allowedEvmVersion, evmVersion) {
		return evmVersion, true
	} else {
		return "byzantium", false
	}
}

func getContractInfo() {

}

func getContracts() {

}

func optimizeValue(value interface{}) string {
	switch rtype := value.(type) {
	case bool:
		if rtype == true{
			return "1"
		} else {
			return "0"
		}
	case string:
		if rtype == "true"{
			return "1"
		} else if rtype == "false" {
			return "0"
		} else {
			panic("")
		}
	default:
		panic("")
	}
}

func optimizationRuns(params map[string]string) string{
	return ""
}

func createSourceFile(code []byte) (string, error) {
	randomID := strconv.Itoa(rand.Int())
	tempDir := os.TempDir() + randomID + NewContractName
	fh, err := os.Create(tempDir)
	if err != nil {
		return "", err
	}
	defer fh.Close()
	_, err = fh.Write(code)
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

func deleteSourceFile(path string) bool {
	err := os.Remove(path)
	if err != nil && os.IsNotExist(err) || err == nil {
		return true
	} else {
		return false
	}
}
package solidity

import (
	//"fmt"
	//"github.com/gatechain/smart_contract_verifier/lib/service/rest"
	contract "github.com/gatechain/smart_contract_verifier/lib/smart_contract"
	//"math/rand"
	//"os"
	//"os/exec"
	"strings"
)

const allowedEvmVersion = "homestead,tangerineWhistle,spuriousDragon,byzantium,constantinople,petersburg,istanbul,default"
const NewContractName = "New.sol"

// Module responsible to compile the Solidity code of a given Smart Contract.
/* raw data
Compiles a code in the solidity command line.

Returns a `Map`.

## Examples

    iex(1)> Explorer.SmartContract.Solidity.CodeCompiler.run([
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
    ...>      optimize: false, evm_version: "byzantium"
    ...>  ])
    {
      :ok,
      %{
        "abi" => [
          %{
            "constant" => false,
            "inputs" => [%{"name" => "x", "type" => "uint256"}],
            "name" => "set",
            "outputs" => [],
            "payable" => false,
            "stateMutability" => "nonpayable",
            "type" => "function"
          },
          %{
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
    }
*/

func Run(params map[string]string) {
	//name := GetValue(params, "name", "")
	//compilerVersion := GetValue(params, "compiler_version", "")
	//code := GetValue(params, "code", "")
	//optimize := GetValue(params, "optimize", "0")
	//optimizationRuns := optimizationRuns(params)
	evmVersion := GetValue(params, "evm_version", LatestEvmVersion())

	checkedEvmVersion, _ := IsEvmVersionAllowed(evmVersion)

	path, err := contract.EnsureExists(checkedEvmVersion)
	if err != nil {
		panic("")
	}
	if path != "" {
		//dir := rest.RegisterApp.GetPath(rest.LocalCompilerRootDir, "priv/compile_solc.js")
		//createSourceFile(code),
		//path,
		//optimizeValue(optimize),
		//optimizationRuns,
		//NewContractName,
		//checkedEvmVersion

		//cmd := exec.Command("node", )
		//cmd.Stdout = os.Stdout
		//output, err := cmd.Output()
		//fmt.Print(output)
		//if err != nil {
		//	panic(err)
		//}
	}
}

func GetValue(dict map[string]string, key, value string) string {
	v, ok := dict[key]
	if ok{
		return v
	} else {
		return value
	}
}

func LatestEvmVersion() string {
	version := AllowedEvmVersions()
	return version[len(version) - 1]
}

func AllowedEvmVersions() []string {
	return strings.Split(allowedEvmVersion, ",")
}

func IsEvmVersionAllowed(evmVersion string) (string, bool) {
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

func createSourceFile(code string) {
	//randonID := rand.Int()
	//tempDir := os.TempDir() + "solidity_source" + string() + ".sol"
}
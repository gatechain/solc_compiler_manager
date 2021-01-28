package cmds

import (
	"fmt"
	"github.com/gatechain/smart_contract_verifier/lib"
	"github.com/gatechain/smart_contract_verifier/lib/compiler"
	"github.com/gatechain/smart_contract_verifier/lib/compiler/solidity"
	"github.com/spf13/cobra"
	"strings"
)

func CompileCMD() *cobra.Command {
	compileCmd := &cobra.Command{
		Use:   "compile",
		Short: "compile solidity source",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get scope
			scopes, _ := cmd.Flags().GetString(solidity.FlagScope)
			err := checkScope(scopes)
			if err != nil {
				return err
			}
			// get name
			name, err := cmd.Flags().GetString(solidity.FlagName)
			if err != nil {
				return err
			}
			// get compile version
			version := args[0]
			executePath, err := compiler.EnsureExists(version)
			if err != nil {
				fmt.Println(err)
				return err
			}
			// get file path
			filePath := args[1]
			if !lib.FileExist(filePath) {
				fmt.Printf("given file not exist, path: %s \n", filePath)
			}
			_, err = solidity.LocalRun(executePath, filePath, name, scopes)
			return err
		},
	}
	compileCmd.Flags().String(
		solidity.FlagScope, "",
		"Choose your abi, hashes, bin compile type for specific output or combined output",
	)
	compileCmd.Flags().String(
		solidity.FlagName, "",
		"Define contract name in your solidity script",
	)

	_ = compileCmd.MarkFlagRequired(solidity.FlagScope)
	_ = compileCmd.MarkFlagRequired(solidity.FlagName)
	return compileCmd
}

func checkScope(scope string) error {
	if scope == "" {
		return fmt.Errorf("missing flag, --scope required")
	}
	scopes := strings.Split(scope, ",")
	for _, s := range scopes {
		if s == "" {
			return fmt.Errorf("scope error: %s", scope)
		}
		if s == "bin" || s == "abi" || s == "hashes" {
			continue
		} else {
			return fmt.Errorf("unknown scope: %s, choose between: abi,bin,hashes", s)
		}
	}
	return nil
}


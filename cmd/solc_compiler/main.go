package main

import (
	"fmt"
	"github.com/gatechain/smart_contract_verifier/lib"
	"github.com/gatechain/smart_contract_verifier/lib/compiler"
	"github.com/gatechain/smart_contract_verifier/lib/compiler/solidity"
	"github.com/spf13/cobra"
	"strings"
)

func main() {
	//Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "solc_compiler",
		Short: "Command line interface for smart contract compiler",
	}

	rootCmd.AddCommand(
		FetchCMD(),
		DeleteCMD(),
		CompileCMD(),
	)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func FetchCMD() *cobra.Command {
	fetchCmd := &cobra.Command{
		Use:   "fetch",
		Short: "fetch solc execute file by given version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := compiler.EnsureExists(args[0])
			if path != "" {
				fmt.Println("Download file: ", path)
			}
			return err
		},
	}
	return fetchCmd
}

func DeleteCMD() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete solc execute file by given version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := compiler.Delete(args[0])
			return err
		},
	}
	return deleteCmd
}

func CompileCMD() *cobra.Command {
	compileCmd := &cobra.Command{
		Use:   "compile",
		Short: "compile solidity source",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get flag
			scopes, _ := cmd.Flags().GetString(solidity.FlagScope)
			err := checkScope(scopes)
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

			solidity.LocalRun(executePath, filePath, scopes)
			return err
		},
	}
	compileCmd.Flags().String(
		solidity.FlagScope, "",
		"Choose your abi, hashes, bin compile type for specific output or combined output",
	)
	_ = compileCmd.MarkFlagRequired(solidity.FlagScope)
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


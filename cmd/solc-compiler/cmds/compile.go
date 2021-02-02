package cmds

import (
	"fmt"
	"github.com/gatechain/solc_compiler_manager/lib"
	"github.com/gatechain/solc_compiler_manager/lib/compiler"
	"github.com/gatechain/solc_compiler_manager/lib/compiler/solidity"
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
			scopes := mustUnpackScope(cmd)
			// get name
			name := mustUnpackName(cmd)
			// get evm version
			evmVersion := mustUnpackEvmVersion(cmd)
			// get optimize
			optimize := mustUnpackOptimize(cmd)
			// get optimize runs
			optimizeRuns := mustUnpackOptimizeRuns(cmd)

			// get compile version
			version := args[0]
			executePath, err := compiler.EnsureExists(version)
			if err != nil {
				return err
			}
			// get file path
			filePath := args[1]
			if !lib.FileExist(filePath) {
				return fmt.Errorf("given file not exist, path: %s \n", filePath)
			}
			_, err = solidity.LocalRun(executePath, filePath, name, scopes, evmVersion, optimize, optimizeRuns)
			return err
		},
	}

	registerFlags(compileCmd)

	return compileCmd
}

func registerFlags(cmd *cobra.Command) {
	_ = cmd.MarkFlagRequired(solidity.FlagScope)
	_ = cmd.MarkFlagRequired(solidity.FlagName)

	cmd.Flags().String(
		solidity.FlagScope, "",
		"Choose your abi, hashes, bin compile options for specific or combined output",
	)
	cmd.Flags().String(
		solidity.FlagName, "",
		"Choose compiled contract in your solidity script",
	)
	cmd.Flags().String(
		solidity.FlagEvmVersion, "default",
		"Choose compiler evm version",
	)
	cmd.Flags().Bool(
		solidity.FlagOptimize, false,
		"Choose compile is optimize or not",
	)
	cmd.Flags().Int(
		solidity.FlagOptimizeRuns, 200,
		"Choose compiler optimize-runs times",
	)
}

func mustUnpackName(cmd *cobra.Command) string {
	// get name
	name, err := cmd.Flags().GetString(solidity.FlagName)
	if err != nil {
		panic(err)
	}
	return name
}

func mustUnpackScope(cmd *cobra.Command) string {
	// get scope
	scopes, err := cmd.Flags().GetString(solidity.FlagScope)
	if err != nil {
		panic(err)
	}
	err = checkScope(scopes)
	if err != nil {
		panic(err)
	}
	return scopes
}

func mustUnpackEvmVersion(cmd *cobra.Command) string {
	// get evm version
	evmVersion, err := cmd.Flags().GetString(solidity.FlagEvmVersion)
	if err != nil {
		panic(err)
	}
	return evmVersion
}

func mustUnpackOptimize(cmd *cobra.Command) bool {
	// get optimize
	optimize, err := cmd.Flags().GetBool(solidity.FlagOptimize)
	if err != nil {
		panic(err)
	}
	return optimize
}

func mustUnpackOptimizeRuns(cmd *cobra.Command) int {
	// get evm version
	optimizeRuns, err := cmd.Flags().GetInt(solidity.FlagOptimizeRuns)
	if err != nil {
		panic(err)
	}
	return optimizeRuns
}

// checkScope check given scope flag format
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
			return fmt.Errorf("unknown scope: %s, choose between: abi,bin,hashes, please seperate it with comma without spaces", s)
		}
	}
	return nil
}


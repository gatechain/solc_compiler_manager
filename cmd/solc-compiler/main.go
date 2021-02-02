package main

import (
	"fmt"
	"github.com/gatechain/solc_compiler_manager/cmd/solc-compiler/cmds"
	"github.com/gatechain/solc_compiler_manager/lib"
	"github.com/spf13/cobra"
)

func main() {
	//Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "solc-compiler",
		Short: "Command line interface for solidity compiler",
	}

	rootCmd.AddCommand(
		cmds.InitCMD(),
		cmds.FetchCMD(),
		cmds.DeleteCMD(),
		cmds.CompileCMD(),
		cmds.ServiceCMD(lib.Cdc),
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}
}

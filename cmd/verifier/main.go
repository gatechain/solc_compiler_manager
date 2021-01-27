package main

import (
	"fmt"
	contract "github.com/gatechain/smart_contract_verifier/lib/smart_contract"
	"github.com/spf13/cobra"
)

func main() {
	//Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "verifier",
		Short: "Command line interface for smart contract verification",
	}

	rootCmd.AddCommand(
		FetchCMD(),
	)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func FetchCMD() *cobra.Command {
	fetchCmd := &cobra.Command{
		Use:   "fetch",
		Short: "fetch solc source with given compiler version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := contract.EnsureExists(args[0])
			if path != "" {
				fmt.Println("Download file: ", path)
			}
			return err
		},
	}
	return fetchCmd
}
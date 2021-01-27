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
		DeleteCMD(),
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
			path, err := contract.EnsureExists(args[0])
			if path != "" {
				fmt.Println("Download file: ", path)
			}
			return err
		},
	}
	return fetchCmd
}

func DeleteCMD() *cobra.Command {
	fetchCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete solc execute file by given version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := contract.Delete(args[0])
			return err
		},
	}
	return fetchCmd
}

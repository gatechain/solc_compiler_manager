package cmds

import (
	"fmt"
	"github.com/gatechain/smart_contract_verifier/lib/compiler"
	"github.com/spf13/cobra"
)

func FetchCMD() *cobra.Command {
	fetchCmd := &cobra.Command{
		Use:   "fetch",
		Short: "fetch solidity compiler by given version",
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

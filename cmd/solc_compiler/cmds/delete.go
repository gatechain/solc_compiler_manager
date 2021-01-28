package cmds

import (
	"github.com/gatechain/smart_contract_verifier/lib/compiler"
	"github.com/spf13/cobra"
)

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

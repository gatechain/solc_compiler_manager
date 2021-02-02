package cmds

import (
	"github.com/gatechain/solc_compiler_manager/lib/compiler"
	"github.com/spf13/cobra"
)

func DeleteCMD() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete solidity compiler by given version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := compiler.Delete(args[0])
			return err
		},
	}
	return deleteCmd
}

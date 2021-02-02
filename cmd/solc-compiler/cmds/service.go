package cmds

import (
	"github.com/gatechain/solc_compiler_manager/lib"
	"github.com/gatechain/solc_compiler_manager/lib/service/rest"
	"github.com/gatechain/solc_compiler_manager/lib/service/rpc"
	"github.com/spf13/cobra"
)

// ServeCommand will start the application REST service.
// It takes a codec to create a RestServer object and a function to register all necessary routes.
func ServiceCMD(cdc *lib.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rs := rest.NewRestServer(cdc)

			rpc.RegisterRoutes(rs)

			// Start the rest server and return error if one exists
			ListenAddr, _ := cmd.Flags().GetString(rest.FlagListenAddr)
			MaxOpenConnections, _ := cmd.Flags().GetInt(rest.FlagMaxOpenConnections)
			RPCReadTimeout, _ := cmd.Flags().GetUint(rest.FlagRPCReadTimeout)
			RPCWriteTimeout, _ := cmd.Flags().GetUint(rest.FlagRPCWriteTimeout)
			UnsafeCORS, _ := cmd.Flags().GetBool(rest.FlagUnsafeCORS)

			err = rs.Start(
				ListenAddr,
				MaxOpenConnections,
				RPCReadTimeout,
				RPCWriteTimeout,
				UnsafeCORS,
			)
			return err
		},
	}
	return rest.RegisterRestServerFlags(cmd)
}

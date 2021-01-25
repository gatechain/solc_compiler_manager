package verifier

import (
)

func main() {
	// Configure cobra to sort commands
	//cobra.EnableCommandSorting = false
	//
	//// Instantiate the codec for the command line application
	//cdc := app.MakeCodec()
	//
	//// TODO: setup keybase, viper object, etc. to be passed into
	//// the below functions and eliminate global vars, like we do
	//// with the cdc
	//
	//rootCmd := &cobra.Command{
	//	Use:   "verifier",
	//	Short: "Command line interface for interacting with gated",
	//}
	//
	//// Add --chain-id to persistent flags and mark it required
	//rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	//rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
	//	return initConfig(rootCmd)
	//}
	//accountCmd := authcmd.AccountCommands(cdc)
	//accountCmd.AddCommand(revocablecmd.PublishMultiSigAccountCmd(cdc))
	//// Construct Root Command
	//rootCmd.AddCommand(
	//	rpc.StatusCommand(),
	//	client.ConfigCmd(app.DefaultCLIHome),
	//	//queryCmd(cdc),
	//	txCmd(cdc),
	//	blockCmd(cdc),
	//	client.LineBreak,
	//	lcd.ServeCommand(cdc, registerRoutes),
	//	client.LineBreak,
	//	revocableTXCmd(cdc),
	//	client.LineBreak,
	//	accountCmd,
	//	vaultAccountCmd(cdc),
	//	client.LineBreak,
	//	tokenCmd(cdc),
	//	tradingPairCmd(cdc),
	//	client.LineBreak,
	//	foundationCmd(cdc),
	//	client.LineBreak,
	//	//proposalCmd(cdc),
	//	stakingCmd(cdc),
	//	slashingCmd(cdc),
	//	distributionCmd(cdc),
	//	validatorCmd(cdc),
	//	client.LineBreak,
	//	dexClient.NewModuleClient(cdc).GetDexCmd(),
	//	client.LineBreak,
	//	version.Cmd,
	//	client.NewCompletionCmd(rootCmd, true),
	//	evmrpc.EvmCommand(cdc),
	//)
	//
	//// Add flags and prefix all env exposed with GA
	//executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)
	//
	//err := executor.Execute()
	//if err != nil {
	//	fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
	//	os.Exit(1)
	//}
}
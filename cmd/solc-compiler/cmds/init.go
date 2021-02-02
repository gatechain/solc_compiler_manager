package cmds

import (
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/gatechain/solc_compiler_manager/lib"
	"github.com/gatechain/solc_compiler_manager/lib/compiler"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"os"
)

// ServeCommand will start the application REST service.
// It takes a codec to create a RestServer object and a function to register all necessary routes.
func InitCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init solidity verifier",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := lib.CompilerLocalHomeDir()
			fetchAll, err := cmd.Flags().GetBool(lib.SolcFetchAll)
			platform, err := cmd.Flags().GetString(lib.LocalPlatForm)
			if err != nil {
				return err
			}
			if !checkPlatform(platform) {
				return fmt.Errorf("please choose platform between %s and %s \n", lib.SolcMacOSX, lib.SolcLinux)
			}
			initConfig := make(lib.LocalConfig)
			initConfig[lib.LocalPlatForm] = platform

			err = writeConfig(home + lib.LocalConfigName, initConfig)
			if err != nil {
				return err
			}

			// init fetch all compiler version on specific platform
			if fetchAll {
				err = compiler.FetchAllVersion(3)  // n for parallel number
				if err != nil {
					return err
				}
			}

			fmt.Println("verifier init success")
			return nil
		},
	}
	cmd.Flags().BoolP(lib.SolcFetchAll, "a", false, "Fetch all compiler version when project init")
	cmd.Flags().String(lib.LocalPlatForm, "", "The platform for the server")
	_ = cmd.MarkFlagRequired(lib.LocalPlatForm)
	return cmd
}

func checkPlatform(platform string) bool {
	if platform == lib.SolcLinux || platform == lib.SolcMacOSX {
		return true
	} else {
		return false
	}
}

// write local config
func writeConfig(path string, params lib.LocalConfig) error {
	if lib.FileExist(path) {
		if inputIsTty() {
			overwrite, err := askOverwrite()
			if err != nil {
				return err
			}
			if overwrite {
				err = lib.WriteJson(path, params)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("file not change")
			}
			return nil
		} else {
			return fmt.Errorf("unsupport method")
		}
	} else {
		err := lib.WriteJson(path, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func askOverwrite() (bool, error) {
	overwrite, err := speakeasy.Ask("init file already exist, overwrite? Y/N")
	if err != nil {
		return false, err
	}
	if overwrite == "Y" {
		return true, nil
	} else if overwrite == "N" {
		return false, nil
	} else {
		fmt.Println("input Y or N")
		return askOverwrite()
	}
}

// inputIsTty returns true if we have an interactive prompt,
// where we can disable echo and request to repeat the password.
// If false, we can optimize for piped input from another command
func inputIsTty() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}
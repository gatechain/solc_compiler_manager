package lib

import (
	"os"
	"os/user"
)

// const for completed urls
const (
	SolcBinApiUrl 			= "https://binaries.soliditylang.org/"
	SolcLinux				= "linux-amd64"
	SolcMacOSX				= "macosx-amd64"
	SolcListVersions		= "list.json"
	SolcVersionLatest		= "latest"
	SolcFetchAll			= "fetch-all"
	LocalCompilerRootDir	= "solc_compilers/"
	LocalSolcCompilerDir 	= "stored/"
	LocalPlatForm			= "platform"
	LocalConfigName			= ".config"
)

// file path
func CompilerLocalStoreDir() string {
	current, _ := user.Current()
	storeDir := current.HomeDir + "/" + LocalCompilerRootDir + LocalSolcCompilerDir
	exist := FileExist(storeDir)
	if !exist {
		err := os.MkdirAll(storeDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return storeDir
}

func CompilerLocalHomeDir() string {
	current, _ := user.Current()
	homeDir := current.HomeDir + "/" + LocalCompilerRootDir
	exist := FileExist(homeDir)
	if !exist {
		err := os.MkdirAll(homeDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return homeDir
}
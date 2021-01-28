package lib

import (
	"os"
	"os/user"
	"syscall"
)

// const for completed urls
const (
	SolcBinApiUrl 			= "https://binaries.soliditylang.org/"
	//SolcPlatform			= "linux-amd64"
	SolcPlatform			= "macosx-amd64"
	SolcListVersions		= "list.json"
	LocalCompilerRootDir	= "solc_compilers/"
	LocalSolcCompilerDir 	= "stored/"
	SolcVersionLatest		= "latest"
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

// file exists
func FileExist(path string) bool {
	err := syscall.Access(path, syscall.F_OK)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type Application map[string]string

func (RegisterApp Application) HaveEnv(env string) bool {
	_, ok := RegisterApp[env]
	return ok
}

func (RegisterApp Application) GetEnv(env string) string {
	return RegisterApp[env]
}

func (RegisterApp Application) GetPath(env, suffix string) string {
	return RegisterApp[env] + RegisterApp[suffix]
}

func (RegisterApp *Application) Set(key, value string) {
	(*RegisterApp)[key] = value
}

// init
var RegisterApp Application

func init() {
	//app := New()
	//
	//// register urls
	//registerPath(&app)
	//
	//RegisterApp = app
}

func registerPath(app *Application) {
	app.Set(SolcBinApiUrl, SolcBinApiUrl)
	app.Set(SolcPlatform, SolcPlatform)
	app.Set(SolcListVersions, SolcListVersions)
	app.Set(LocalCompilerRootDir, LocalCompilerRootDir)
	app.Set(LocalSolcCompilerDir, LocalSolcCompilerDir)
}



package client

type Application map[string]string

func New() Application {
	return make(Application)
}

func (RegisterApp Application) GetEnv(path string) string {
	return RegisterApp[path]
}

func (RegisterApp Application) HavePath(path string) bool {
	_, ok := RegisterApp[path]
	return ok
}

// init
var RegisterApp Application

func init() {
	app := New()

	// register urls
	app["explorer/solc_bin_api_url"] = "https://solc-bin.ethereum.org"

	RegisterApp = app
}



package smart_contract

import (
	"fmt"
	"github.com/gatechain/smart_contract_verifier/lib"
	"github.com/gatechain/smart_contract_verifier/lib/service/rest"
	"net/url"
	"os"
	"os/exec"
	"time"
)

func EnsureExists(version string) (string, error){
	path := filePath(version)

	if lib.FileExist(path) && version != lib.SolcVersionLatest {
		return path, nil
	} else {
		versions := fetchVersions()
		if release := versionReleased(versions, version); release != ""{
			HandleCall(version, release)
			return path, nil
		} else {
			return "", fmt.Errorf("given version not exists")
		}
	}
}

func HandleCall(version, release string) {
	path := filePath(version)
	if needFetch(version, path) {
		success := download(release)
		if !success {
			panic(fmt.Sprintf("create file %s filed", path))
		}
	}
}

func filePath(version string) string {
	return lib.CompilerLocalStoreDir() + fmt.Sprintf("solc-%s-v%s", lib.SolcPlatform, version)
}

// fetch version
func needFetch(version, path string) bool {
	switch version {
	case "latest":
		state, err := os.Stat(path)
		if err != nil {
			return true
		}
		return time.Now().Sub(state.ModTime()) > 30
	default:
		return !lib.FileExist(path)
	}
}

// version interactions
func client() rest.Client {
	checkedUrl, err := url.Parse(lib.SolcBinApiUrl)
	if err != nil {
		panic(fmt.Sprintf("base url format illegal, get: %s", lib.SolcBinApiUrl))
	}
	var client = rest.MakeRestClient(*checkedUrl, "")
	return client
}

func fetchVersions() rest.SolcVersion {
	client := client()
	var resp rest.SolcVersion
	err := client.FetchVersions(&resp)
	if err != nil {
		return rest.SolcVersion{}
	}
	return resp
}

func versionReleased(versions rest.SolcVersion, version string) string {
	// TODO need more cache
	if version == lib.SolcListVersions {
		return versions.Releases[versions.LatestRelease]
	}
	for _, build := range versions.Builds {
		if build.Version == version {
			return build.Path
		}
	}
	return ""
}

func download(version string) bool {
	client := client()
	return client.Download(version)
}

func Delete(version string) error {
	path := filePath(version)
	if lib.FileExist(path) {
		return deleteVersion(path)
	} else {
		fmt.Printf("version: %s not exist", version)
		return nil
	}
}

func deleteVersion(path string) error {
	fmt.Printf("delete file: %s \n", path)
	cmd := fmt.Sprintf("rm %s", path)
	command := exec.Command("bash", "-c", cmd)
	err := command.Run()
	if err != nil {
		fmt.Printf("delete file filed: %s \n", path)
		return err
	} else {
		print("delete success")
		return nil
	}
}
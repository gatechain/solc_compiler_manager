package compiler

import (
	"fmt"
	"github.com/gatechain/smart_contract_verifier/lib"
	"github.com/gatechain/smart_contract_verifier/lib/service/rest"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"
)

func EnsureExists(version string) (string, error){
	if err := lib.CheckVersionFormat(version); err != nil {
		return "", fmt.Errorf("version foramt not match, need version like: v0.8.0, got %s", version)
	}

	path := lib.FilePath(version)

	if lib.FileExist(path) && version != lib.SolcVersionLatest {
		return path, nil
	} else {
		versions := fetchVersions()
		if release := versionReleased(versions, version); release != ""{
			HandleCall(version, release)
			return path, nil
		} else {
			return "", fmt.Errorf("no release can match given version")
		}
	}
}

func HandleCall(version, release string) {
	path := lib.FilePath(version)
	if needFetch(version, path) {
		success := download(release)
		if !success {
			panic(fmt.Sprintf("download file %s filed", path))
		}
	}
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

func fetchVersions() lib.SolcVersion {
	client := client()
	var resp lib.SolcVersion
	err := client.FetchVersions(&resp)
	if err != nil {
		return lib.SolcVersion{}
	}
	return resp
}

func versionReleased(versions lib.SolcVersion, version string) string {
	// TODO need more cache
	if version == lib.SolcVersionLatest {
		return versions.Releases[versions.LatestRelease]
	}
	for _, build := range versions.Builds {
		if "v" + build.Version == version {
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
	path := lib.FilePath(version)
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

func FetchAllVersion() error {
	versions := fetchVersions()
	var wg sync.WaitGroup
	wg.Add(len(versions.Releases))
	for version, release := range versions.Releases {
		go goHandleCall(&wg, version, release)
	}
	wg.Wait()
	return nil
}

func goHandleCall(wg *sync.WaitGroup, version, release string) {
	HandleCall(version, release)
	wg.Done()
}


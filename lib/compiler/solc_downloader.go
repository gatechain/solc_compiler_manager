package compiler

import (
	"fmt"
	"github.com/gatechain/solc_compiler_manager/lib"
	"github.com/gatechain/solc_compiler_manager/lib/service/rest"
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
			HandleCall(version, release, nil)
			return path, nil
		} else {
			return "", fmt.Errorf("no release can match given version")
		}
	}
}

func HandleCall(version, release string, bar *lib.Bar) {
	path := lib.FilePath(version)
	if needFetch(version, path) {
		success := download(release, bar)
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

func download(version string, bar *lib.Bar) bool {
	client := client()
	return client.Download(version, bar)
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

func FetchAllVersion(n ...int) error {
	versions := fetchVersions()
	var wg sync.WaitGroup
	wg.Add(len(versions.Releases))
	plot := lib.NewMultiProgressBar("Start Download: ")

	var jobNum int
	if len(n) > 0 {
		jobNum = n[0]
	} else {
		jobNum = len(versions.Releases)
	}
	jobs := make(chan string, jobNum)
	for version, release := range versions.Releases {
		// skip downloaded file
		path := lib.FilePath(version)
		if lib.FileExist(path) && version != lib.SolcVersionLatest {
			continue
		}

		jobs <- version
		bar := plot.NewBar(version + "\t", 100)
		go goHandleCall(&jobs, &wg, bar, version, release)
	}
	wg.Wait()
	return nil
}

func goHandleCall(jobs *chan string, wg *sync.WaitGroup, bar *lib.Bar, version, release string) {
	HandleCall(version, release, bar)
	<-*jobs
	defer wg.Done()
}


package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

// version format
var (
	regexVersion   	= `v[0-9]+\.[0-9]+\.[0-9]+`
	regexCommitSep	= `\+{1}`
	regexCommit   	= `commit\.[a-z0-9]+`
	VersionMatch 	= regexp.MustCompile(fmt.Sprintf(`^%s$`, regexVersion))
	LongVersionMatch 	= regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, regexVersion, regexCommitSep, regexCommit))
)

func CheckVersionFormat(version string) error {
	matches := VersionMatch.FindStringSubmatch(version)
	if matches == nil || len(matches) != 1{
		return fmt.Errorf("invalid version '%s'", version)
	}
	return  nil
}

func CheckLongVersionFormat(version string) (string, string, error) {
	matches := LongVersionMatch.FindStringSubmatch(version)
	if matches == nil || len(matches) != 3{
		return "", "", fmt.Errorf("invalid version '%s'", version)
	}
	return matches[1], matches[2], nil
}

func CheckVersionCommit(version string, commit string) bool {
	path := FilePath(version)
	// readlink returns real file name
	cmd := fmt.Sprintf("readlink %s", path)
	command := exec.Command("bash", "-c", cmd)
	output, err := command.Output()
	if err != nil {
		return false
	}
	localFile := strings.Trim(filepath.Base(string(output)), "\n")
	seps := strings.Split(localFile, "+")
	if len(seps) != 2 {
		return false
	}
	return commit == seps[1]
}

// return local stored compiler file path
func FilePath(version string) string {
	return CompilerLocalStoreDir() + fmt.Sprintf("solc-%s-%s", GetPlatform(), version)
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

// get os platform
func GetPlatform() string {
	config := make(LocalConfig)
	home := CompilerLocalHomeDir()
	err := ReadJson(home + LocalConfigName, config)
	if err != nil {
		return ""
	}
	return config[LocalPlatForm]
}

// read local file in json format
func ReadJson(path string, o interface{}) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()
	bz, err := ioutil.ReadAll(fh)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bz, &o)
	if err != nil {
		return err
	}
	return nil
}


// write local file in json format
func WriteJson(path string, o interface{}) error {
	//var fh * os.File
	var err error
	bz, err := json.Marshal(o)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bz, 0644)
	//fh, err = os.OpenFile(path, os.O_APPEND, os.ModePerm)
	if err != nil {
		//if os.IsNotExist(err) {
		//	fh, err = os.Create(path)
		//	if err != nil {
		//		return err
		//	}
		//}
		return err
	}
	//defer fh.Close()
	//_, err = fh.Write(bz)
	//if err != nil {
	//	return err
	//}
	return nil
}
package compiler

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFetchVersions(t *testing.T) {
	bz, err:= json.MarshalIndent(FetchVersions(), "", "\t")
	if err != nil {
		t.Errorf("fetch versions error, %s", err)
	}
	fmt.Println(string(bz))
}
package lib

import (
	"fmt"
	"regexp"
)

// version format
var (
	regexVersion   	= `v[0-9]+\.[0-9]+\.[0-9]+`
	regexCommitSep	= `\+{1}`
	regexCommit   	= `commit\.[a-z]|[0-9]+`
	VersionMatch 	= regexp.MustCompile(fmt.Sprintf(`^%s$`, regexVersion))
	LongVersionMatch 	= regexp.MustCompile(fmt.Sprintf(`^%s%s%s$`, regexVersion, regexCommitSep, regexCommit))
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
	if matches == nil || len(matches) != 2{
		return "", "", fmt.Errorf("invalid version '%s'", version)
	}
	return matches[0], matches[1], nil
}

func CheckVersionCommit(version string, commit string) bool {
	return true
}
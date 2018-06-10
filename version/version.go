package version

import "fmt"

var Commit string

func VersionInfo() string {
	return fmt.Sprintf("version: %s", Commit)
}

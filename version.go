package hookworm

import (
	"fmt"
	"os"
	"path"
)

var (
	VersionString string
	progName      string
)

func init() {
	progName = path.Base(os.Args[0])
}

func printVersion() {
	fmt.Println(progVersion())
}

func progVersion() string {
	if VersionString == "" {
		VersionString = "<unknown>"
	}

	return fmt.Sprintf("%s %s", progName, VersionString)
}

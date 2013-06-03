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
	if VersionString == "" {
		VersionString = "<unknown>"
	}

	fmt.Printf("%s %s\n", progName, VersionString)
}

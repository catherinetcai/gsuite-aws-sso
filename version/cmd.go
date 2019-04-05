package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd returns the version of the tool
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the tool",
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}

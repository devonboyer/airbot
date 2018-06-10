package app

import (
	"fmt"

	"github.com/devonboyer/airbot/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.VersionInfo())
	},
}

package cmd

import (
	"fmt"

	"github.com/notedownorg/notedown/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display detailed version information including build details.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.GetInfo().String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

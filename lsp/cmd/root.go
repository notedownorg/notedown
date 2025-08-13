package cmd

import (
	"fmt"
	"os"

	"github.com/notedownorg/notedown/lsp/pkg/constants"
	"github.com/notedownorg/notedown/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   constants.LanguageServerName,
	Short: "Notedown Language Server Protocol implementation",
	Long: `A Language Server Protocol (LSP) implementation for Notedown flavored Markdown.
Provides language features like completion, diagnostics, and navigation for Notedown documents.`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Println(version.GetInfo().String())
			return
		}
		fmt.Println("Notedown LSP Server")
		fmt.Println("Use --help for available commands")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
}

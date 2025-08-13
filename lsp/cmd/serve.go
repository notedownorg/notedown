package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/version"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the LSP server",
	Long: `Start the Notedown Language Server Protocol server.
The server communicates via stdin/stdout using the LSP protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Notedown LSP server...")
		
		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(os.Stdout)
		
		mux := lsp.NewMux(reader, writer, version.Get())
		if err := mux.Run(); err != nil {
			log.Fatalf("LSP server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	
	serveCmd.Flags().StringP("log-file", "l", "", "Path to log file (default: stderr)")
	serveCmd.Flags().StringP("log-level", "", "info", "Log level (debug, info, warn, error)")
}
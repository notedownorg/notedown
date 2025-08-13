package cmd

import (
	"bufio"
	"os"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/lsp/pkg/notedownls"
	"github.com/notedownorg/notedown/pkg/log"
	"github.com/notedownorg/notedown/pkg/version"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the LSP server",
	Long: `Start the Notedown Language Server Protocol server.
The server communicates via stdin/stdout using the LSP protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFile, _ := cmd.Flags().GetString("log-file")
		logFormat, _ := cmd.Flags().GetString("log-format")

		level := log.ParseLevel(logLevel)
		format := log.ParseFormat(logFormat)
		var logger *log.Logger
		var err error

		if logFile != "" {
			logger, err = log.NewFile(logFile, level, format)
			if err != nil {
				panic(err)
			}
		} else {
			logger = log.NewLsp(level, format)
		}

		logger.WithScope("lsp/cmd").Info("starting notedown lsp server", "version", version.Get())

		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(os.Stdout)

		// Create Notedown-specific LSP server
		server := notedownls.NewServer(version.Get(), logger)

		// Create mux and set the server
		mux := lsp.NewMux(reader, writer, version.Get(), logger)
		mux.SetServer(server)

		if err := mux.Run(); err != nil {
			logger.WithScope("lsp/cmd").Error("lsp server failed", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("log-file", "l", "", "Path to log file (default: stderr)")
	serveCmd.Flags().StringP("log-level", "", "info", "Log level (debug, info, warn, error)")
	serveCmd.Flags().StringP("log-format", "", "text", "Log format (text, json)")
}

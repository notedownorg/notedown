// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"strings"

	"github.com/notedownorg/notedown/cmd/initialize"
	"github.com/spf13/cobra"
)

var (
	Version    string
	CommitHash string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "notedown",
	Short: "Tools for note-taking and organization",
}

func init() {
	rootCmd.AddCommand(initialize.RootCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func version() string {
	var b strings.Builder

	if Version == "" {
		b.WriteString("dev")
	} else {
		b.WriteString(Version)
	}

	if CommitHash != "" {
		b.WriteString("-")
		b.WriteString(CommitHash)
	}

	return b.String()
}

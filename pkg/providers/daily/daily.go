// Copyright 2024 Notedown Authors
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

package daily

import (
	"log/slog"
	"path/filepath"
	"strings"
	"time"
)

const MetadataKey = "daily"

type identifier struct {
	path    string
	version string
}

// By default we will set line to -1 to default to end of file
func NewIdentifier(path string, version string) identifier {
	return identifier{path: path, version: version}
}

func (i identifier) String() string {
	// Pipe separators are good enough for now but may need to be changed as pipes
	// are technically valid (although unlikely to actually be used) in unix file paths
	// We may want to consider an actual encoding scheme for this in the future.
	var builder strings.Builder
	builder.WriteString(i.path)
	builder.WriteString("|")
	builder.WriteString(i.version)
	return builder.String()
}

type Daily struct {
	name       string
	identifier identifier
	date       time.Time
}

func NewDaily(identifier identifier) Daily {
	name := strings.TrimSuffix(filepath.Base(identifier.path), filepath.Ext(identifier.path))

	// TODO: Support more than just YYYY-MM-DD
	date, err := time.Parse("2006-01-02", name)
	if err != nil {
		slog.Error("failed to parse date from daily note name", "name", name, "identifier", identifier, "error", err)
	}

	return Daily{
		identifier: identifier,
		name:       name,
		date:       date,
	}
}

func (d Daily) Identifier() identifier {
	return d.identifier
}

func (d Daily) Name() string {
	return d.name
}

func (d Daily) Path() string {
	return d.identifier.path
}

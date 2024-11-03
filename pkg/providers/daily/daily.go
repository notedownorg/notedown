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
	"path/filepath"
	"strings"
)

type Identifier struct {
	path    string
	version string
}

// By default we will set line to -1 to default to end of file
func NewIdentifier(path string, version string) Identifier {
	return Identifier{path: path, version: version}
}

func (i Identifier) String() string {
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
	identifier Identifier
}

func NewDaily(identifier Identifier) Daily {
	return Daily{
		identifier: identifier,
		name:       strings.TrimSuffix(filepath.Base(identifier.path), filepath.Ext(identifier.path)),
	}
}

func (d Daily) Identifier() Identifier {
	return d.identifier
}

func (d Daily) Name() string {
	return d.name
}

func (d Daily) Path() string {
	return d.identifier.path
}

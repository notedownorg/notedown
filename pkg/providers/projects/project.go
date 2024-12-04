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

package projects

import (
	"path/filepath"
	"strings"
)

const MetadataKey = "project"

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

const (
    StatusKey = "status"
    NameKey   = "name"
)

type Status string

const (
	Active    Status = "active"
	Archived  Status = "archived"
	Abandoned Status = "abandoned"
	Blocked   Status = "blocked"
	Backlog   Status = "backlog"
)

var statusMap = map[string]Status{
	"active":    Active,
	"archived":  Archived,
	"abandoned": Abandoned,
	"blocked":   Blocked,
	"backlog":   Backlog,
}

type Project struct {
	name       string
	identifier identifier
	status     Status
}

type ProjectOption func(*Project)

func NewProject(identifier identifier, opts ...ProjectOption) Project {
	name := strings.TrimSuffix(filepath.Base(identifier.path), filepath.Ext(identifier.path))
	p := Project{
		identifier: identifier,
		name:       name,
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func NewProjectFromProject(project Project, opts ...ProjectOption) Project {
	p := Project{
		identifier: project.identifier,
		name:       project.name,
		status:     project.status,
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func WithStatus(status Status) ProjectOption {
	return func(p *Project) {
		p.status = status
	}
}

// Encapsulate name changes, external consumers should use client.RenameProject
func withName(name string) ProjectOption {
	return func(p *Project) {
		p.name = name
	}
}

func (p Project) Identifier() identifier {
	return p.identifier
}

func (p Project) Name() string {
	return p.name
}

func (p Project) Path() string {
	return p.identifier.path
}

func (p Project) Status() Status {
	return p.status
}

func (p Project) String() string {
	return p.Name()
}

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

package writer

import (
	"fmt"
	"path/filepath"

	"github.com/notedownorg/notedown/pkg/configuration"
)

type Client struct {
	root   string
	config *configuration.WorkspaceConfiguration
}

func NewClient(ws *configuration.Workspace) (*Client, error) {
	root := configuration.ExpandPath(ws.Location)
	config, err := configuration.EnsureWorkspaceConfiguration(root)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure workspace configuration: %w", err)
	}
	return &Client{root: root, config: config}, nil
}

func (c Client) abs(doc string) string {
	return filepath.Join(c.root, doc)
}

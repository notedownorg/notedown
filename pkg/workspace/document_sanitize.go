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

package workspace

import (
	"github.com/notedownorg/notedown/pkg/configuration"
)

// Should be run prior to saving a document to ensure that the document is in a valid state
func Sanitize(config configuration.WorkspaceConfiguration, d *Document) {
	if d.Metadata == nil || len(d.Metadata) == 0 {
		return
	}
	sanitizeTags(d.Metadata, config.Tags.DefaultFormat)
}

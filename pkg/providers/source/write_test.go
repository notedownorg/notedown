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

package source

import (
	"testing"

	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/workspace"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	client, _ := buildClient(loadEvents(),
		test.Validators{
			Create: []test.CreateValidator{
				func(doc workspace.Document) error {
					assert.Equal(t, "library/source.md", doc.Path())
					assert.Equal(t, workspace.Metadata{
						workspace.MetadataTypeKey: MetadataKey,
						TitleKey:                  "source",
						FormatKey:                 Article,
						UrlKey:                    "example.com",
					}, doc.Metadata)
					return nil
				},
			},
			Delete: []test.DeleteValidator{
				func(path string) error {
					assert.Equal(t, "library/source.md", path)
					return nil
				},
			},
		},
	)
	assert.NoError(t, client.CreateSource("library/source.md", "source", Article, "example.com"))
	assert.NoError(t, client.DeleteSource(Source{path: "library/source.md"}))
}

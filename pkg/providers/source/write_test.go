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

package source_test

import (
	"testing"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/providers/source"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	client, _ := buildClient(loadEvents(),
		test.Validators{
			Create: []test.CreateValidator{
				func(doc writer.Document, metadata reader.Metadata, content []byte, feed chan reader.Event) error {
					assert.Equal(t, writer.Document{Path: "library/source.md"}, doc)
					assert.Equal(t, reader.Metadata{
						reader.MetadataTypeKey: source.MetadataKey,
						source.TitleKey:        "source",
						source.FormatKey:       source.Article,
						source.UrlKey:          "example.com",
					}, metadata)
					assert.Equal(t, []byte("# source\n\n"), content)
					return nil
				},
			},
			Delete: []test.DeleteValidator{
				func(doc writer.Document) error {
					assert.Equal(t, writer.Document{Path: "library/source.md"}, doc)
					return nil
				},
			},
		},
	)
	assert.NoError(t, client.CreateSource("library/source.md", "source", source.Article, "example.com"))
	assert.NoError(t, client.DeleteSource(source.NewArticle(source.NewIdentifier("library/source.md", ""), "source", "example.com")))
}

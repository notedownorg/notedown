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

package reader

import (
	"fmt"

	"github.com/a-h/parse"
	"sigs.k8s.io/yaml"
)

type Metadata map[string]interface{}

func (m Metadata) Type() string {
	if m == nil {
		return ""
	}
	typeValue, ok := m[MetadataType]
	if !ok {
		return ""
	}
	if res, ok := typeValue.(string); ok {
		return res
	}
	return ""
}

const (
	MetadataType = "type"
)

type Document struct {
	Metadata Metadata
	Contents []byte
	Checksum string

	// Last updated is used to determine if we even need to bother reading the file from disk
	// It should only be used internally and shouldn't be exposed to the consumer
	lastUpdated int64
}

var parseDocument = func() func(string) (Document, error) {
	return func(input string) (Document, error) {
		p := parse.NewInput(input)
		res, ok, err := documentParser().Parse(p)
		if err != nil {
			return Document{}, fmt.Errorf("unable to parse document: %w", err)
		}
		if !ok {
			return Document{}, fmt.Errorf("unable to parse document")
		}

		return res, nil
	}
}

var documentParser = func() parse.Parser[Document] {
	return parse.Func(func(in *parse.Input) (Document, bool, error) {
		var res Document

		// Look for frontmatter
		frontmatterTuple, ok, err := parse.SequenceOf2(parse.AtLeast(0, parse.Whitespace), parseFrontmatter).Parse(in)
		if err != nil {
			return Document{}, false, err
		}
		if ok {
			err := yaml.Unmarshal(frontmatterTuple.B, &res.Metadata)
			if err != nil {
				return Document{}, false, fmt.Errorf("unable to parse frontmatter: %w", err)
			}
		}

		// Parse the rest of the document
		rem, _, err := parse.StringUntil(parse.EOF[string]()).Parse(in)
		if err != nil {
			return Document{}, false, fmt.Errorf("unable to parse contents: %w", err)
		}
		res.Contents = []byte(rem)

		return res, true, nil
	})
}

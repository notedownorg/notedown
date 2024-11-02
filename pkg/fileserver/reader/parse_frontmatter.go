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
	"encoding/json"
	"fmt"

	"github.com/a-h/parse"
	"github.com/notedownorg/notedown/pkg/parsers"
	"sigs.k8s.io/yaml"
)

type frontmatter []byte

var frontMatterKeyword = parse.String("---")

var frontMatterOpen = parse.StringFrom(frontMatterKeyword, parsers.RemainingInlineWhitespace, parse.StringFrom(parse.AtLeast(1, parse.NewLine)))
var frontMatterClose = parse.StringFrom(parse.StringFrom(parse.AtLeast(0, parse.NewLine)), frontMatterKeyword, parsers.RemainingInlineWhitespace)

var parseFrontmatter parse.Parser[frontmatter] = parse.Func(func(in *parse.Input) (frontmatter, bool, error) {
	// Read and discard the front matter open.
	if _, ok, err := frontMatterOpen.Parse(in); err != nil || !ok {
		return nil, false, err
	}

	// Read up to the front matter close.
	contents, _, err := parse.StringUntil(frontMatterClose).Parse(in)
	if err != nil {
		return nil, false, err
	}

	// Technically, the front matter could be empty...
	// If it isnt empty, we need to check that it is valid yaml.
	if len(contents) != 0 {
		// To do this we need to convert it to json and then use the stdlib to check it.
		jsn, err := yaml.YAMLToJSON([]byte(contents))
		if err != nil {
			return nil, false, fmt.Errorf("couldnt validate frontmatter yaml: %w", err)
		}
		if !json.Valid(jsn) {
			return nil, false, fmt.Errorf("front matter is not valid yaml")
		}
	}

	// Read and discard the front matter close
	if _, ok, err := frontMatterClose.Parse(in); err != nil || !ok {
		return nil, false, err
	}

	// Discard final newline if it exists
	parse.StringFrom(parse.AtMost(1, parse.NewLine)).Parse(in)

	return frontmatter(contents), true, nil
})

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

package parsers

import (
	"fmt"
	"time"

	"github.com/a-h/parse"
	"github.com/notedownorg/notedown/pkg/ast"
	"sigs.k8s.io/yaml"
)

var Document = func(path string, version string, relativeTo time.Time) func(string) (ast.Document, error) {
	return func(input string) (ast.Document, error) {
		p := parse.NewInput(input)
		res, ok, err := DocumentParser(path, version, relativeTo).Parse(p)
		if err != nil {
			return ast.Document{}, fmt.Errorf("unable to parse document: %w", err)
		}
		if !ok {
			return ast.Document{}, fmt.Errorf("unable to parse document")
		}

		return res, nil
	}
}

var DocumentParser = func(path, version string, relativeTo time.Time) parse.Parser[ast.Document] {
	return parse.Func(func(in *parse.Input) (ast.Document, bool, error) {
		var res ast.Document

		// Look for frontmatter
		frontmatterTuple, ok, err := parse.SequenceOf2(parse.AtLeast(0, parse.Whitespace), Frontmatter).Parse(in)
		if err != nil {
			return ast.Document{}, false, err
		}
		if ok {
			err := yaml.Unmarshal(frontmatterTuple.B, &res.Metadata)
			if err != nil {
				return ast.Document{}, false, fmt.Errorf("unable to parse frontmatter: %w", err)
			}
			res.Markers.ContentStart = in.Position().Line + 1
		}

		// Parse the rest of the file looking for blocks
		blocks, ok, err := parse.Until(Block(path, version, relativeTo), parse.EOF[string]()).Parse(in)
		if err != nil {
			return ast.Document{}, false, err
		}
		for _, b := range blocks {
			res.Tasks = append(res.Tasks, b.Tasks...)
		}

		return res, true, nil
	})
}

type block struct {
	Tasks []ast.Task
}

var Block = func(path, version string, relativeTo time.Time) parse.Parser[block] {
	return parse.Func(func(in *parse.Input) (block, bool, error) {
		var res block

		// Drop any leading newline
		_, _, err := parse.NewLine.Parse(in)

		for {
			task, ok, err := Task(path, version, relativeTo).Parse(in)
			if err != nil {
				return block{}, false, err
			}
			if !ok {
				break
			}
			res.Tasks = append(res.Tasks, task)

		}

		// Process the input until the next newline or EOF as the current line isnt a task
		_, _, err = parse.StringUntil(newLineOrEOF).Parse(in)
		if err != nil {
			return block{}, false, err
		}

		return res, true, nil
	})
}

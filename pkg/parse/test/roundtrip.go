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

package test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/liamawhite/parse/core"
	. "github.com/liamawhite/parse/test"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	"github.com/stretchr/testify/assert"
)

func VerifyRoundTrip(t *testing.T, tests []ParserTest[[]ast.Block]) {
	for _, test := range tests {
		t.Run(fmt.Sprintf("Roundtrip: %v", test.Name), func(t *testing.T) {
			in := NewInput(test.Input)
			blocks, _, _ := test.Parser(in)

			var got strings.Builder
			for _, block := range blocks {
				got.WriteString(block.Markdown())
				got.WriteString("\n")
			}
			assert.Equal(t, test.Input, got.String(), "Roundtrip mismatch from original input")

			// Do multiple round trips to ensure that the AST is stable
			// Only run if we haven't already failed to avoid spamming the output
			if !t.Failed() {
				for i := 0; i < 10; i++ {
					in := NewInput(got.String())
					blocks, _, _ := test.Parser(in)

					got.Reset()
					for _, block := range blocks {
						got.WriteString(block.Markdown())
						got.WriteString("\n")
					}
					assert.Equal(t, test.Input, got.String(), "Mismatch after multiple roundtrips")
				}
			}
		})

	}

}

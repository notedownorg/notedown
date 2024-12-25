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

package blocks

import (
	"crypto/md5"
	"io"

	"github.com/notedownorg/notedown/pkg/parse/ast"
)

// tracker is used to track the state of a block to verify if it has been modified.
type tracker struct {
	checkpoint string
}

func newTracker(block ast.Block) *tracker {
	return &tracker{
		checkpoint: hash(block),
	}
}

func hash(b ast.Block) string {
	h := md5.New()
	io.WriteString(h, b.Markdown())
	return string(h.Sum(nil))
}

func (t *tracker) Modified(b ast.Block) bool {
	return t.checkpoint != hash(b)
}

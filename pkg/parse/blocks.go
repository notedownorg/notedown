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

package parse

import (
	"time"

	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	"github.com/notedownorg/notedown/pkg/parse/blocks"
	"github.com/notedownorg/notedown/pkg/parse/traversal"
)

func Blocks(relativeTo time.Time) Parser[[]ast.Block] {
	return func(in Input) ([]ast.Block, bool, error) {
		return traversal.DepthFirstParser(blocks.TopLevelParser(blocks.NewContext(relativeTo)), EOF[ast.Block]())(in)
	}
}

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

package workspace

import (
	"fmt"
	"time"

	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse"
	"github.com/notedownorg/notedown/pkg/parse/blocks"
)

var documentParser = func(relativeTo time.Time) Parser[Document] {
	return func(in Input) (Document, bool, error) {
		var res Document

		// Parse the document
		blks, ok, err := parse.Blocks(relativeTo)(in)
		if err != nil {
			return Document{}, false, fmt.Errorf("unable to parse blocks: %w", err)
		}
		if !ok {
			return Document{}, false, fmt.Errorf("unable to parse blocks: no blocks found")
		}

		// If the first block is a frontmatter block hoist the metadata
		if len(blks) > 0 {
			if frontmatter, ok := blocks.Frontmatter(blks[0]); ok {
				res.Metadata = frontmatter.Metadata()
				blks = blks[1:]
			}
		}

		res.Blocks = blks

		return res, true, nil
	}
}

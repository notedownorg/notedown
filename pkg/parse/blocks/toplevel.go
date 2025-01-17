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
	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
)

func TopLevelParser(ctx Context) []Parser[ast.Block] {
	return []Parser[ast.Block]{
		FrontmatterParser,
		ThematicBreakParser,

		ListParser(ctx),
		BlockQuoteParser(ctx),
		CodeBlockIndentedParser,
		CodeBlockFencedParser,

		HeadingAtxParser(ctx, 1, 1),
		HeadingAtxParser(ctx, 2, 2, HeadingAtxParser(ctx, 1, 1)),
		HeadingAtxParser(ctx, 3, 3, HeadingAtxParser(ctx, 1, 2)),
		HeadingAtxParser(ctx, 4, 4, HeadingAtxParser(ctx, 1, 3)),
		HeadingAtxParser(ctx, 5, 5, HeadingAtxParser(ctx, 1, 4)),
		HeadingAtxParser(ctx, 6, 6, HeadingAtxParser(ctx, 1, 5)),

		HeadingSetextParser(ctx, 1, 1),
		HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),

		HTMLParser,

		BlankLineParser,
		ParagraphParser(ctx),
	}
}

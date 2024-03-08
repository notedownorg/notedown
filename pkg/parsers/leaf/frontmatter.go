package leaf

import (
	"encoding/json"
	"fmt"

	"github.com/a-h/parse"
	"sigs.k8s.io/yaml"
)

type FrontMatter []byte

var frontMatterKeyword = parse.String("---")

var frontMatterOpen = parse.StringFrom(frontMatterKeyword, parse.StringFrom(parse.AtLeast(0, parse.RuneIn(" \t"))), parse.StringFrom(parse.AtLeast(1, parse.NewLine)))
var frontMatterClose = parse.StringFrom(parse.StringFrom(parse.AtLeast(0, parse.NewLine)), frontMatterOpen)

var FrontMatterParser parse.Parser[FrontMatter] = parse.Func(func(in *parse.Input) (FrontMatter, bool, error) {
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

	// Read and discard the front matter close.
	if _, ok, err := frontMatterClose.Parse(in); err != nil || !ok {
		return nil, false, err
	}

	return FrontMatter(contents), true, nil
})

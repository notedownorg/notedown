package parsers

import (
	"fmt"
	"time"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/api"
	"sigs.k8s.io/yaml"
)

var Document = func(relativeTo time.Time) func (string) (api.Document, error) {
    return func(input string) (api.Document, error) {
        p := parse.NewInput(input)
        res, ok, err := DocumentParser(relativeTo).Parse(p)
        if err != nil {
            return api.Document{}, fmt.Errorf("unable to parse document: %w", err)
        }
        if !ok {
            return api.Document{}, fmt.Errorf("unable to parse document")
        }
        return res, nil
    }
}

var DocumentParser = func(relativeTo time.Time) parse.Parser[api.Document] {
	return parse.Func(func(in *parse.Input) (api.Document, bool, error) {
		var res api.Document

		// Look for frontmatter
		frontmatterTuple, ok, err := parse.SequenceOf2(parse.AtLeast(0,parse.Whitespace), Frontmatter).Parse(in)
		if err != nil {
			return api.Document{}, false, err
		}
		if ok {
			err := yaml.Unmarshal(frontmatterTuple.B, &res.Metadata)
			if err != nil {
				return api.Document{}, false, fmt.Errorf("unable to parse frontmatter: %w", err)
			}
		}

        // Parse the rest of the file looking for blocks
        blocks, ok, err := parse.Until(Block(relativeTo), parse.EOF[string]()).Parse(in)
        if err != nil {
            return api.Document{}, false, err
        }
        for _, b := range blocks {
            res.Tasks = append(res.Tasks, b.Tasks...)
        }

		return res, true, nil
	})
}

type block struct {
	Tasks []api.Task
}

var Block = func(relativeTo time.Time) parse.Parser[block] {
	return parse.Func(func(in *parse.Input) (block, bool, error) {
		var res block

        // Drop any leading newline
        _, _, err := parse.NewLine.Parse(in)

		// TODO: do something more correct than blindly looking for tasks in the input
        // Read until we find a task
        _, _, err = parse.StringUntil(Task(relativeTo)).Parse(in)
        if err != nil {
            return block{}, false, err
        }

        task, ok, err := Task(relativeTo).Parse(in)
        if err != nil {
            return block{}, false, err
        }
        if ok {
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

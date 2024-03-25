package parsers

import (
	"fmt"
	"time"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/api"
	"sigs.k8s.io/yaml"
)

var Document = func(relativeTo time.Time) func(string) (api.Document, error) {
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
		frontmatterTuple, ok, err := parse.SequenceOf2(parse.AtLeast(0, parse.Whitespace), Frontmatter).Parse(in)
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
            if len(b.Tasks) > 0 {
                fmt.Printf("adding block with %v tasks\n", len(b.Tasks))
                res.Tasks = append(res.Tasks, b.Tasks...)
            }
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

        parent := api.Task{Name: "root", Indent: -1}
        // stack := []api.Task{parent}
        // previousTask := api.Task{Name: "foo", Indent: 0}
        for {
            task, ok, err := Task(relativeTo).Parse(in)
            if err != nil {
                return block{}, false, err
            }
            if !ok {
                fmt.Println("saw something not a task, breaking out of block")
                break
            }
            if !parent.AddChild(&task) {
                fmt.Printf("failed to add task %q %v somehow", task.Name, task.Indent)
            }
            // if task.Indent == previousTask.Indent {
            //     parent.SubTasks = append(parent.SubTasks, task)
            //     fmt.Printf("same indent: appending %q to %q sub tasks (%v)\n", task.Name, parent.Name, len(parent.SubTasks))
            //     previousTask = task
            //     continue
            // }
            // if task.Indent > previousTask.Indent {
            //     fmt.Printf("increased indent: pushing %q to stack and adding to parent %q\n", task.Name, parent.Name)
            //     stack = append(stack, previousTask)
            //     parent = previousTask
            //     parent.SubTasks = append(parent.SubTasks, task)
            //     previousTask = task
            //     continue
            // }
            // fmt.Printf("indent decreased. parent %v (%v -> %v)\n", parent.Indent, previousTask.Indent, task.Indent)
            // if task.Indent == parent.Indent {
            //     // we're back to the level of the parent
            //     fmt.Printf("back to parent indent level\n")
            //     parent, stack = stack[len(stack)-1], stack[:len(stack)-1]
            //     parent.SubTasks = append(parent.SubTasks, task)
            //     continue
            // }
            // // indent decreased below parent
            // for {
            //     parent, stack = stack[len(stack)-1], stack[:len(stack)-1]
            //     fmt.Printf("new parent %q %v\n", parent.Name, parent.Indent)
            //     if task.Indent >= parent.Indent {
            //         parent.SubTasks = append(parent.SubTasks, task)
            //         break
            //     }
            // }


        }

        // for i := len(stack) - 1; i >= 1; i-- {
        //     stack[i-1].SubTasks = append(stack[i-1].SubTasks, stack[i])
        // }
        root := parent
        if len(root.SubTasks) > 0 {
            fmt.Printf("%q %v\n", root.Name, len(root.SubTasks))
            for _, t := range root.SubTasks {
                fmt.Printf("\t%q\n", t.Name)
                res.Tasks = append(res.Tasks, *t)
            }
            // res.Tasks = append(res.Tasks, stack[0].SubTasks...)
        }

        // Process the input until the next newline or EOF as the current line isnt a task
        _, _, err = parse.StringUntil(newLineOrEOF).Parse(in)
        if err != nil {
            return block{}, false, err
        }
    
        fmt.Printf("end of block with %v tasks\n", len(res.Tasks))
		return res, true, nil
	})
}

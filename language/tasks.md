# Task Lists

Task lists provide a structured way to manage tasks and to-do items in Notedown Flavored Markdown. Tasks are rendered as interactive checkboxes and can be nested within hierarchical list structures.

## Basic Syntax

Task lists use GitHub Flavored Markdown-style checkbox syntax within list items:

```markdown
- [ ] Unchecked task
- [x] Checked task
- [X] Checked task (uppercase also supported)
```

## List Context Requirements

Task checkboxes are only recognized within list items. They must appear:

1. At the beginning of a list item text
2. Within unordered (`-`, `*`, `+`) or ordered (`1.`, `2.`, etc.) lists
3. As the first element of the list item content

### Valid Contexts

```markdown
- [x] In unordered list
* [ ] In unordered list (asterisk)
+ [x] In unordered list (plus)

1. [ ] In ordered list
2. [x] In ordered list
```

## Nested Tasks

Tasks can be nested within hierarchical list structures, this is considered a subtask:

```markdown
- [x] Main project task
  - [x] Completed subtask
  - [ ] Pending subtask
    - [ ] Nested subtask
    - [x] Another nested task
- [ ] Another main task
  1. [x] Ordered subtask
  2. [ ] Another ordered subtask
```

## Mixed Lists

Task items can be mixed with regular list items:

```markdown
- [x] Task item
- Regular list item
- [ ] Another task
- Another regular item
```

## Integration with Wikilinks

Tasks can contain wikilinks for cross-referencing:

```markdown
- [x] Implement [[features/completion-engine|completion engine]]
- [ ] Add [[features/definition-provider|go-to-definition]]
- [ ] Update [[docs/api-reference|API reference]]
- [ ] Test with [[test-data/large-workspace|large workspaces]]
```

## Custom Task States

Notedown supports customizable task states beyond the standard `[ ]` (todo) and `[x]` (done) through workspace configuration. This allows teams to define task states that match their workflow.

### Workspace Configuration

Create a `.notedown/settings.yaml` (or `.notedown/settings.json`) file in your workspace root to define custom task states:

```yaml
tasks:
  states:
    - value: " "           # What goes inside [brackets]
      name: "todo"         # Human-readable name
      aliases: []          # Optional aliases
    - value: "x"
      name: "done"
      aliases: ["X", "✓", "✔"]
    - value: "/"
      name: "in-progress"
      aliases: ["wip", "WIP"]
    - value: "-"
      name: "cancelled"
      aliases: ["~", "cancelled"]
    - value: "?"
      name: "question"
      aliases: ["Q"]
    - value: "important"
      name: "high-priority"
      aliases: ["!", "urgent"]
```

### Configuration Rules

1. **Unique Values**: All `value` and `aliases` must be unique across all states
2. **No Reserved Characters**: Values and aliases cannot contain `]`
3. **Required Fields**: `value` and `name` are required for each state
4. **Optional Fields**: `aliases` are optional
5. **Format Support**: Both YAML (`.yaml`) and JSON (`.json`) formats are supported, with YAML preferred

### Custom State Examples

With the configuration above, you can use any of these task state syntaxes:

```markdown
- [ ] Standard todo task
- [x] Completed task
- [/] Work in progress task
- [wip] Alternative in-progress syntax
- [-] Cancelled task
- [?] Question or needs clarification
- [important] High priority task
- [!] Alternative high priority syntax
```

### Configuration Discovery

The configuration file is discovered by walking up the directory tree from the current file, looking for a `.notedown/` directory containing `settings.yaml` or `settings.json`. This allows:

- Project-specific task state definitions
- Shared team workflows through version control
- Fallback to default states when no configuration exists

### Default Configuration

When no workspace configuration is found, Notedown uses these default task states:

```yaml
tasks:
  states:
    - value: " "
      name: "todo"
    - value: "x" 
      name: "done"
      aliases: ["X"]
```

This ensures backward compatibility with standard GitHub Flavored Markdown task lists.

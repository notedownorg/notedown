# Workspace Configuration

Notedown supports workspace-level configuration through `.notedown/settings.yaml` or `.notedown/settings.json` files. This allows teams to customize behavior and define shared conventions.

## Configuration Discovery

Configuration files are discovered by searching up the directory tree from the current file location, looking for a `.notedown/` directory containing either:

1. `settings.yaml` (preferred)
2. `settings.json` (fallback)

If both files exist, YAML takes precedence. The search stops at the first `.notedown/` directory found, establishing the workspace root.

## File Structure

The `.notedown/` directory should be placed at your workspace root (typically your project's repository root) and can contain:

```
workspace-root/
├── .notedown/
│   ├── settings.yaml    # Main configuration (YAML preferred)
│   └── settings.json    # Alternative format (JSON)
├── docs/
├── src/
└── README.md
```

## Configuration Schema

### Complete Example (YAML)

```yaml
tasks:
  states:
    - value: " "
      name: "todo"
      aliases: []
    - value: "x"
      name: "done"
      conceal: "✅"
      aliases: ["X", "✓", "✔"]
    - value: "/"
      name: "in-progress"
      conceal: "⏳"
      aliases: ["wip", "WIP", "working"]
    - value: "-"
      name: "cancelled"
      conceal: "❌"
      aliases: ["~", "cancelled", "skip"]
    - value: "?"
      name: "question"
      conceal: "❓"
      aliases: ["Q", "unclear"]
    - value: "waiting"
      name: "blocked"
      conceal: "⏸️"
      aliases: ["blocked", "hold"]
    - value: "!"
      name: "urgent"
      conceal: "🔥"
      aliases: ["urgent", "priority", "asap"]
```

### Complete Example (JSON)

```json
{
  "tasks": {
    "states": [
      {
        "value": " ",
        "name": "todo",
        "aliases": []
      },
      {
        "value": "x",
        "name": "done",
        "conceal": "✅",
        "aliases": ["X", "✓", "✔"]
      },
      {
        "value": "/",
        "name": "in-progress", 
        "conceal": "⏳",
        "aliases": ["wip", "WIP", "working"]
      },
      {
        "value": "-",
        "name": "cancelled",
        "conceal": "❌",
        "aliases": ["~", "cancelled", "skip"]
      }
    ]
  }
}
```

## Task State Configuration

### Required Fields

- `value`: The text that appears inside `[brackets]` in the markdown
- `name`: Human-readable name for the task state

### Optional Fields

- `conceal`: Visual replacement text for compatible editors (emoji, symbol, or text)
- `aliases`: Array of alternative values that map to the same state

### Validation Rules

1. **Uniqueness**: All `value` and `aliases` entries must be unique across all states
2. **Reserved Characters**: Values and aliases cannot contain `]` (interferes with bracket syntax)
3. **Non-Empty**: Values, names, and aliases cannot be empty strings
4. **Length Limits**: Reasonable length limits for readability (values should be concise)

### Examples of Valid States

```yaml
# Single character states
- value: "x"
  name: "done"

# Multi-character states  
- value: "in-progress"
  name: "work-in-progress"

# States with concealment
- value: "!"
  name: "urgent"
  conceal: "🔥"

# States with aliases
- value: "x"
  name: "done"
  aliases: ["X", "✓", "✔", "complete"]

# Unicode values
- value: "✓"
  name: "checked"
  
# Word-based values
- value: "waiting"
  name: "blocked"
  conceal: "⏸️"
```

## Usage in Markdown

Once configured, task states can be used in any list context:

```markdown
- [x] Completed task
- [ ] Todo task  
- [/] Work in progress
- [wip] Alternative in-progress (via alias)
- [-] Cancelled task
- [?] Question or unclear requirement
- [waiting] Blocked on external dependency
- [!] Urgent priority task
```

## Editor Integration

### Concealment Support

When `conceal` is specified, compatible editors can replace the bracket syntax with the conceal text:

- Raw markdown: `- [x] Task completed`
- Concealed display: `- ✅ Task completed`

### Syntax Highlighting

Editors can use the configuration to provide appropriate syntax highlighting and autocomplete for defined task states.

## Version Control

The `.notedown/settings.yaml` file should be committed to version control to share task state definitions across the team. This ensures consistent task management workflows.

## Migration and Compatibility

### From Default States

Existing markdown with `[ ]` and `[x]` will continue to work without configuration changes. The default configuration includes these standard states.

### Adding Custom States

Custom states can be added incrementally. Existing task lists will continue to function while new states become available for future use.

### Removing States

When removing a task state from configuration:
1. Existing uses in markdown will no longer be recognized as task states
2. They will render as plain text: `[removed-state]`
3. Consider migration scripts for large codebases

## Best Practices

### State Design

1. **Keep Values Concise**: Short values are easier to type and read
2. **Use Meaningful Names**: Names should clearly indicate the task state
3. **Provide Useful Aliases**: Include common variations and shortcuts
4. **Consider Visual Appeal**: Choose appropriate conceal characters

### Team Adoption

1. **Document Your States**: Include examples in project documentation
2. **Start Simple**: Begin with a few essential states, expand gradually
3. **Get Team Buy-in**: Ensure all team members understand the workflow
4. **Regular Review**: Periodically review and refine state definitions

### Common Patterns

```yaml
# Minimal workflow
tasks:
  states:
    - value: " "
      name: "todo"
    - value: "x" 
      name: "done"
      aliases: ["X"]

# Agile workflow
tasks:
  states:
    - value: " "
      name: "backlog"
    - value: "/"
      name: "in-progress"
      conceal: "⏳"
    - value: "r"
      name: "review"
      conceal: "👁️"
    - value: "x"
      name: "done"
      conceal: "✅"

# Detailed workflow
tasks:
  states:
    - value: " "
      name: "todo"
    - value: "/"
      name: "in-progress"
      conceal: "⏳"
    - value: "r"
      name: "review"
      conceal: "👁️"
    - value: "b"
      name: "blocked"
      conceal: "🚫"
    - value: "x"
      name: "done"
      conceal: "✅"
    - value: "-"
      name: "cancelled"
      conceal: "❌"
```
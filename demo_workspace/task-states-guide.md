# Task States Guide

This demo workspace showcases the enhanced task state system with descriptions, aliases, and visual conceal characters.

## Available Task States

### Standard States
- `[ ]` **todo** - A task that needs to be completed
- `[x]` **done** - A completed task ✅
- `[/]` **in-progress** - A task currently being worked on ⏳
- `[?]` **question** - A task that needs clarification or more information ❓
- `[!]` **urgent** - A high priority task that should be completed soon 🔥
- `[-]` **cancelled** - A task that has been cancelled and will not be completed ❌

### Aliases
Each task state supports multiple aliases for flexibility:

**Done aliases**: `X`, `✓`, `✔`, `completed`, `finished`
- `[X] Task completed using uppercase X`
- `[completed] Task marked as completed`
- `[finished] Task that is finished`

**In-progress aliases**: `wip`, `WIP`, `working`
- `[wip] Work in progress task`
- `[WIP] Work in progress (uppercase)`
- `[working] Currently working on this`

**Question aliases**: `Q`, `unclear`, `help-needed`, `blocked`
- `[Q] Quick question marker`
- `[unclear] Something unclear about this task`
- `[blocked] Task is blocked by something`

**Urgent aliases**: `urgent`, `priority`, `critical`, `important`
- `[priority] High priority task`
- `[critical] Critical task`
- `[important] Important to complete`

**Cancelled aliases**: `~`, `cancelled`, `wont-do`, `rejected`
- `[~] Quick cancellation marker`
- `[wont-do] Task we won't do`
- `[rejected] Task was rejected`

## LSP Completion

When using an LSP-enabled editor, you'll get intelligent completion when typing `- [`:

1. **Rich Details**: Each completion shows the task state name and description
2. **Visual Indicators**: Conceal characters (emojis) are shown where available
3. **Alias Information**: Aliases clearly indicate which main state they reference
4. **Smart Filtering**: Type part of a state or alias to filter results

### Example Completion Display:
```
- [   <-- cursor here, trigger completion
  ┌─────────────────────────────────────────────────┐
  │ x     Task state: done - A completed task (✅)  │
  │ X     Task state: done - ... [alias for 'x']   │
  │ /     Task state: in-progress - A task curr...  │
  │ wip   Task state: in-progress - ... [alias]    │
  │ ?     Task state: question - A task that ne...  │
  │ !     Task state: urgent - A high priority...   │
  └─────────────────────────────────────────────────┘
```

## Configuration

The task states are configured in `.notedown/settings.yaml`:

```yaml
tasks:
  states:
    - value: "x"
      name: "done"
      description: "A completed task"
      conceal: "✅"
      aliases: ["X", "completed", "finished"]
```

## Benefits

1. **Consistency**: Standard vocabulary across your workspace
2. **Flexibility**: Multiple ways to express the same state
3. **Clarity**: Descriptions explain what each state means
4. **Visual Appeal**: Conceal characters provide visual distinction
5. **Editor Support**: LSP integration provides intelligent completion

## Examples in Practice

See [[tasks.md]] for extensive examples of these task states in use across a complex project structure with nested tasks and wikilink integration.

Try editing [[test-completion.md]] to experience the completion system yourself!
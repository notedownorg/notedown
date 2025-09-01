# Wikilink Completion

Demonstrates Notedown's intelligent wikilink completion system with live LSP-powered suggestions as you type.

## Demo

![Wikilink Completion Demo](./demo.gif)

## How It Works

When typing `[[` to start a wikilink, Notedown's LSP provides intelligent completions based on:

1. **Existing Files** (highest priority): Files that actually exist in the workspace
2. **Referenced Targets** (medium priority): Targets that are referenced by other wikilinks but don't exist yet
3. **Directory Paths** (lowest priority): Path-based completions for hierarchical organization

## Usage

1. Start typing a wikilink with `[[`
2. Begin typing the target name (e.g., `proj`)
3. Use `Ctrl+X Ctrl+O` to trigger LSP completions 
4. Navigate suggestions with `Ctrl+N` and `Ctrl+P`
5. Press `Ctrl+Y` to accept the completion
6. Complete with `]]`

## Example

In a workspace with files:
- `project-alpha.md`
- `meeting-notes.md`
- `docs/api.md`

Typing `[[pro` would suggest:
- `project-alpha` (existing file match)
- `project-beta` (if referenced elsewhere)
- `projects/` (directory completion)

## Features Demonstrated

- **Smart Filtering**: Completions match based on prefix
- **Prioritization**: Existing files appear first in suggestions
- **Context Awareness**: Suggestions adapt to workspace content
- **Performance**: Fast completion response even in large workspaces
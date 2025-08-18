# Demo Workspace

This is a sample workspace demonstrating Notedown Flavored Markdown (NFM) features and LSP functionality.

## Structure

- `projects/` - Project documentation and notes
- `notes/` - General notes and ideas  
- `docs/` - Documentation files
- `ambiguous-example.md` - Demonstrates diagnostic warnings for ambiguous wikilinks

## Features Demonstrated

- [[Wikilinks]] for internal linking
- Standard Markdown formatting
- Task lists and organization
- Code blocks and syntax highlighting
- **Ambiguous link detection** - Try opening `ambiguous-example.md` to see diagnostic warnings

## Diagnostic Examples

The workspace includes intentional conflicts to demonstrate LSP diagnostics:

- `[[config]]` is ambiguous (could refer to `/config.md` or `/docs/config.md`)
- Other wikilinks like `[[README]]`, `[[docs/architecture]]` are unambiguous

## Development Usage

Run `make dev` to install the language server and open this workspace in Neovim for testing.
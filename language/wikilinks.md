# Internal Links/Wikilinks

Wikilinks provide a simple way to create internal links between documents using double bracket syntax.

## Basic Syntax

```markdown
[[Page Name]]
```

Creates a link to another page in your workspace. The link text displays as "Page Name" and links to a document with that name.

## Features

### Simple Links
```markdown
[[getting-started]]
[[project-overview]]
```

### Links with Display Text
```markdown
[[document-name|Display Text]]
```

Links to `document-name.md` but shows "Display Text" as the link text.

## Behavior

- **Automatic Resolution**: Searches for matching files in your workspace
- **Markdown Only**: Links only resolve to `.md` files
- **Fuzzy Matching**: Intelligent matching of document names

## Examples

```markdown
# My Document

This references [[another-document]] in the workspace.

You can also link to [[projects/project-overview]] or use 
[[getting-started|Getting Started Guide]] for better readability.

All paths are resolved from the workspace root.
```

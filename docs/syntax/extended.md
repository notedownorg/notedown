# Extended Syntax

Extended syntax covers and formatting that is not considered part of CommonMark spec. Some of these features may be correctly interpreted by Markdown tooling as the syntax is borrowed from other Markdown dialects (Obsidian and Jekyll to name a few) but as they are not part of the recognised CommonMark spec it enitrely depends. 

## Frontmatter

```
---
tags: 
  - project/foo
some: string
list:
  - foo
  - bar
number: 1
checkbox: true
---
```

Frontmatter is YAML enclosed by three `-` symbols. It is used to attach structured data to a document, usually (but not exclusively) by programs that interact with your workspace.

Comments are not supported and will not be persisted if a file is modified.

Excluding comments, any valid YAML is considered part of the specification, but nesting is usually avoided by convention.

## Tags

Tags are only supported in Frontmatter (for now!) under the tags key and follow [Obsidian's hierarchical tags](https://help.obsidian.md/Editing+and+formatting/Tags#Nested+tags) format and syntax.

Tags must:

- Contain at least one letter.
- Contain only alphanumeric characters, underscores (`_`), hyphens (`-`) or forward slash (`/`) for nested/hierarchical tags.

As spaces are prohibited in tags you can use either `kebab-case` (default), `snake_case`, `camelCase` or `PascalCase`. Tools should use the value set in Notedown workspace configuration.

### "Types" via tags

Notedown tooling uses hierarchical tags to identify notes that they should care about. These are identified using the first level in the hierarchy e.g. sources use the `source/title` tag. Using tags in this way allows multiple tools to manage a single note in the context they care about it and allows for polymorhpism via duck/structural typing. This design was heavily influenced by No Boilerplate's philosophy described [here](https://youtu.be/B0yAy2j-9V0).

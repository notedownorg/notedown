# Extended Syntax

Extended syntax covers and formatting that is not considered part of CommonMark spec. Some of these features may be correctly interpreted by Markdown tooling as the syntax is borrowed from other Markdown dialects (Obsidian and Jekyll to name a few) but as they are not part of the recognised CommonMark spec it enitrely depends. 

## Frontmatter

```
---
type: project
some: string
list:
  - foo
  - bar
number: 1
checkbox: true
---
```

Frontmatter is YAML enclosed by three `-` symbols. It is used to attach structured data to a document, usually (but not exclusively) by programs that interact with your workspace.

Comments are not currently supported. If you would like to put comments in your frontmatter please open an issue.

By convention, the `type` keyword (`type: project` in the example) is used by programs interacting with your workspace to determine whether they need to process a given document or they can safely ignore it.

Excluding comments, any valid YAML is considered part of the specification, but nesting is usually avoided by convention.

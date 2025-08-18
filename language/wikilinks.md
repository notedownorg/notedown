# Wikilinks (Internal Links)

Wikilinks provide a simple and powerful way to create internal links between documents in your workspace using double bracket syntax. They enable seamless navigation and cross-referencing within your Notedown knowledge base.

## Basic Syntax

Wikilinks use double bracket notation to create internal document references:

```markdown
[[document-name]]
[[document-name|Display Text]]
```

### Simple Links

```markdown
[[getting-started]]
[[project-overview]]
[[api-reference]]
```

Creates links that display the target name and link to the corresponding document.

### Links with Custom Display Text

```markdown
[[getting-started|Getting Started Guide]]
[[api-reference|API Documentation]]
[[project-overview|Project Overview]]
```

Uses the pipe character (`|`) to separate the target from the display text.

## Syntax Rules

### Valid Wikilink Format

- Must start with `[[` and end with `]]`
- Target name cannot be empty
- Target name cannot contain `]` or `|` characters
- Optional display text after pipe separator (`|`)
- Whitespace around target and display text is automatically trimmed

### Character Restrictions

**Valid characters in targets:**
- Letters, numbers, hyphens, underscores
- Forward slashes for directory paths
- Spaces (converted to hyphens for file matching)

**Invalid syntax:**
```markdown
[[]] # Empty target
[[target|]] # Empty display text
[[tar]get]] # Brackets in target
[[target|dis|play]] # Multiple pipes
[[[target]]] # Triple brackets
```

## Target Resolution

### File Matching

Wikilinks resolve to Markdown files (`.md`) in your workspace:

```markdown
[[project-alpha]] → project-alpha.md
[[meeting-notes]] → meeting-notes.md
[[user-endpoints]] → api/user-endpoints.md
[[api/user-endpoints]] → api/user-endpoints.md
```

### Resolution Strategy

1. **Exact filename match**: `[[user-guide]]` → `user-guide.md`
2. **Path-based match**: `[[docs/api-reference]]` → `docs/api-reference.md`
3. **Fuzzy matching**: Intelligent matching of similar names
4. **Case-insensitive**: `[[User-Guide]]` matches `user-guide.md`

### Directory Traversal

Wikilinks support directory paths for organized workspaces:

```markdown
[[docs/architecture]]      # Links to docs/architecture.md
[[projects/alpha/readme]]  # Links to projects/alpha/readme.md
```

## Semantic Meaning

In Notedown Flavored Markdown, wikilinks represent:

- **Structural relationships** between documents
- **Knowledge connections** and cross-references
- **Navigation pathways** through content
- **Dependency tracking** between components


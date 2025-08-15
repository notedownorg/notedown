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

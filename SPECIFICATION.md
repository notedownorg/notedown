# (WIP) Specification

`nl` (NoteLang) is a markdown derivative specifically adapted for use with Personal Knowledge Management. The goal of the language is to give your note taking an IDE-like experience. As such, it is optimised for use in programming IDEs/text editors however `nl` workspaces can also be opened as [Obsidian](https://obsidian.md) vaults.

This specification is a work in progress as features are added to the language. For now, assume that all Markdown is valid and any extensions that are not part of the [CommonMark specification](https://spec.commonmark.org/0.31.2/) will be detailed below. Until we reach 1.0, there are no guarantees that all CommonMark will be valid but unless you are doing something uncommon (pun not intented) it will likely work.

## Frontmatter

Frontmatter is a 1-1 mapping with [Obsidian properties](https://help.obsidian.md/Editing+and+formatting/Properties). Certain fields in frontmatter con be configured to allow `nl` to infer structure, for example `type: project`, these are optional. 

## Wikilinks

Wikilinks are a form of internal links (links to notes in the opened workspace/vault) in the form of `[[note name]]`.

## Tasks

Tasks can be created using the following format:
```
- [ ] Write the nl specification
```
Tasks have fields that can be configured by appending to the end of the task. Once a single field has been added, the rest of the line is considered to be field configuration and therefore is not part of the task name. Fields are an unordered list in the format `<key>:<value>`. The task name and values cannot contain `:`.

### Status

Status is indicated by the value enclosed between the square brackets.

```md
- [ ] todo
- [b] blocked
- [/] doing
- [x] done
- [a] abandoned
```

### Fields

Due Dates -> `due:2024-01-01 or d:2024-01-01`
Scheduled Dates -> `scheduled:2024-01-01 or s:2024-01-01`
Recurrences -> `every:2 weeks or e:2 weeks`
Priority -> `priority:1 or p:1`

For example, creating a task to repeat daily starting on the 25th May 2024 with a priority of two, note that the comma separation is optional and.
```
- [ ] Some task due:2024-05-25 every:day priority:2
```

### Project Tasks

Tasks that exist in notes with the frontmatter `type: project` are considered to belong to that project.



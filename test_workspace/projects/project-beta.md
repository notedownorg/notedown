# Project Beta

Secondary project with cross-references to [[project-alpha]].

## Description

This project builds upon the work done in [[project-alpha|Project Alpha]] and extends it with additional functionality.

## Features

1. **Enhanced parsing** - Improved markdown processing
2. **Better integration** - Seamless connection with [[project-alpha]]
3. **Advanced wikilinks** - Support for `[[target|display]]` syntax

## Implementation Notes

```typescript
interface WikiLink {
    target: string;
    display?: string;
}
```

See [[notes/technical-decisions]] for rationale behind design choices.
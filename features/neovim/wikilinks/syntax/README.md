# Wikilink Syntax

This test demonstrates comprehensive wikilink syntax support in Notedown, combining both display text formatting and path-based resolution features.

## Features Demonstrated

### Display Text Syntax (`[[target|display]]`)

- **Format**: `[[target|display]]` where target is used for navigation and display is shown to users
- **Completion**: Works on the target portion (before the `|`) 
- **Navigation**: Go-to-definition uses the target regardless of cursor position
- **Use Case**: Provides readable link text while maintaining precise target resolution

### Path Resolution (`[[path/to/file]]`)

- **Directory Navigation**: Supports nested directory structures like `[[docs/api]]`
- **File Creation**: Automatically creates missing directories and files
- **Extension Handling**: Automatically appends `.md` extension
- **Multi-level Paths**: Handles deep directory structures

### Combined Features

- **Path + Display**: `[[path/to/file|Custom Display]]` combines both features
- **Completion Support**: Directory-aware completion for path-based targets
- **Consistent Navigation**: Same go-to-definition behavior across all syntax variations

## Test Workflow

1. **Existing Path Navigation**: Tests navigation to existing directory-based wikilinks
2. **Display Text Creation**: Demonstrates creating wikilinks with custom display text  
3. **Completion Testing**: Shows completion working on target portion with display text
4. **Path Completion**: Tests directory-based completion for new path wikilinks
5. **Navigation Behavior**: Verifies consistent navigation from both display and target portions
6. **File Creation**: Tests automatic directory and file creation for new paths

## Workspace Structure

```
workspace/
├── index.md                 # Main test file with various wikilink examples
├── docs/
│   └── api.md              # Documentation file for path testing
├── projects/
│   └── notedown.md         # Project file for nested path testing
├── team/
│   └── members.md          # Team directory for completion testing
├── meeting-notes.md        # File for display text target testing  
└── project-alpha.md        # File for display text target testing
```

## Expected Behavior

- ✅ Display text shown in editor but target used for navigation
- ✅ Completion suggestions based on target portion only
- ✅ Path-based wikilinks resolve to correct directory structure
- ✅ Missing directories created automatically
- ✅ Go-to-definition works from any cursor position within wikilink
- ✅ Combined syntax (`[[path|display]]`) works seamlessly
- ⚠️ Wikilink concealment (showing only display text) is enabled but may not be fully visible in VHS recordings due to LSP initialization timing

## Note on Concealment

The Notedown plugin includes wikilink concealment functionality that hides the `[[target|` and `]]` syntax, showing only the display text for better readability. This feature:

- Requires LSP server initialization and conceal range calculation
- Is enabled with `conceallevel = 2` and `concealcursor = "nc"`
- Works reliably in manual testing but may not consistently appear in automated VHS recordings
- Provides a cleaner editing experience by showing `My Project` instead of `[[project-alpha|My Project]]`

This comprehensive test validates that Notedown's wikilink system handles both human-readable display text and flexible path-based organization effectively.
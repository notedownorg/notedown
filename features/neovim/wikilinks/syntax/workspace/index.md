# Wikilink Syntax Test Workspace

This workspace demonstrates comprehensive wikilink syntax including display text and path resolution.

## Display Text Syntax

Examples of wikilinks using the `[[target|display]]` syntax:

- [[project-alpha|My Amazing Project]] - Target: `project-alpha`, Display: "My Amazing Project"
- [[meeting-notes|Team Meeting 2024-01-15]] - Target: `meeting-notes`, Display: "Team Meeting 2024-01-15"  
- [[docs/api|API Reference Guide]] - Target: `docs/api`, Display: "API Reference Guide"

## Directory-based Wikilinks

Test path-based wikilink resolution and navigation:

- [[docs/api]] - Navigate to documentation in `docs/` directory
- [[projects/notedown]] - Navigate to project file in `projects/` directory  
- [[team/members]] - Navigate to team directory structure

## Multi-level Paths with Display Text

Combining both features:

- [[projects/notedown|Notedown Project]] - Path with display text
- [[team/members|Team Directory]] - Directory with custom display

## Features Demonstrated

1. **Display Text**: `[[target|display]]` format for readable link text
2. **Path Resolution**: Directory-based navigation with `[[path/to/file]]`
3. **Completion Support**: Works on target portion, supports directories  
4. **Navigation**: Go-to-definition uses target for both syntax types
5. **File Creation**: Missing directories and files created automatically

## Testing Instructions

1. Test go-to-definition on existing wikilinks (both syntax types)
2. Create new wikilinks with completion support
3. Test navigation from both display text and target portions
4. Verify directory structure creation for non-existent paths
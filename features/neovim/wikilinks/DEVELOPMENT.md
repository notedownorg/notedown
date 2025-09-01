# Wikilinks Development Guide

This guide covers the implementation details and development practices for wikilink feature testing.

## Architecture Overview

Wikilinks in Notedown are implemented through multiple LSP components:

### Core Components

1. **Wikilink Index** (`indexes/wikilink.go`)
   - Maintains registry of all wikilink targets
   - Tracks file existence and reference counts
   - Detects ambiguous targets with multiple matches
   - Provides completion suggestions

2. **LSP Handlers** (`handlers_textDocument.go`)
   - `handleCompletion`: Provides intelligent wikilink completions
   - `handleDefinition`: Enables go-to-definition navigation
   - `handleCodeAction`: Offers quick fixes for ambiguous wikilinks

3. **Diagnostics System** (`diagnostics.go`)
   - Detects broken wikilinks (missing targets)
   - Identifies ambiguous wikilinks (multiple matches)
   - Provides real-time error highlighting

4. **Parser Integration** (`pkg/parser/extensions/wikilink.go`)
   - AST-based wikilink extraction
   - Support for `[[target|display]]` syntax
   - Position tracking for LSP features

## Testing Strategy

### Feature Coverage

Each feature test demonstrates a specific aspect of wikilink functionality:

| Feature | LSP Handler | Index Method | Key Behavior |
|---------|-------------|--------------|--------------|
| basic-completion | `handleCompletion` | `GetTargetsByPrefix` | Show completion suggestions |
| goto-definition | `handleDefinition` | `findFileForTarget` | Navigate to existing/new files |
| file-creation | `handleDefinition` | `createMarkdownFile` | Auto-create missing targets |
| display-text-syntax | Parser | `extractWikilinksFromAST` | Handle pipe separator |
| path-resolution | `resolveTargetPath` | `targetExistsInWorkspace` | Directory-based links |
| ambiguous-detection | Diagnostics | `GetAmbiguousTargets` | Multiple file matches |
| code-actions | `handleCodeAction` | `generateCodeActionsForTarget` | Quick fix suggestions |
| diagnostics | `generateWikilinkDiagnostics` | `GetNonExistentTargets` | Error highlighting |

### Test Workspace Design

Each feature uses carefully crafted workspaces to demonstrate specific scenarios:

- **Existing Files**: Show completion and navigation to real targets
- **Missing Files**: Demonstrate file creation behavior  
- **Ambiguous Targets**: Create scenarios with multiple matching files
- **Directory Structure**: Test path-based wikilinks
- **Display Text**: Include examples of custom display syntax

## VHS Template Patterns

### Common Template Structure

```vhs
# Feature-specific setup
Output "{{.OutputFile}}"
Set FontSize 12
Set Width 1200
Set Height 800
Set TypingSpeed 50ms

Hide
# Setup workspace
Type "cd '{{.WorkspaceDir}}'"
Enter
Sleep 1s

# Start Neovim with plugin
Type "nvim --clean -u '{{.ConfigFile}}' index.md"
Enter  
Sleep 4s
Show

# Feature demonstration
Type "[[target"
Sleep 2s
# Show completions appearing
Ctrl+Space
Sleep 3s

Hide
Type ":qa!"
Enter
```

### Key Template Variables

- `{{.WorkspaceDir}}`: Test workspace path
- `{{.ConfigFile}}`: Neovim configuration with Notedown plugin
- `{{.OutputFile}}`: ASCII output path for golden testing
- `{{.LSPBinary}}`: Built language server binary path

## Development Workflow

### Adding New Wikilink Features

1. **Identify LSP Component**: Determine which handler/method needs testing
2. **Design Test Scenario**: Create workspace files that demonstrate the feature
3. **Write VHS Template**: Show the feature in action with realistic usage
4. **Create Documentation**: Explain what the feature does and how to use it
5. **Validate Coverage**: Ensure the test exercises the target LSP code

### Testing Best Practices

1. **Realistic Scenarios**: Use practical examples developers would encounter
2. **Clear Demonstrations**: Show both successful and edge case behaviors
3. **Proper Timing**: Allow sufficient sleep time for LSP responses
4. **Visual Clarity**: Use appropriate font sizes and screen dimensions
5. **Error Handling**: Test both success and failure paths

### Debugging Tests

Common issues and solutions:

- **LSP Not Starting**: Check binary path and plugin installation
- **Completions Not Showing**: Verify trigger characters and timing
- **Navigation Failures**: Ensure workspace files exist as expected
- **Template Errors**: Validate all template variables are available

## Implementation Details

### Completion System

The completion system provides three tiers of suggestions:

1. **Existing Files** (highest priority): Files that actually exist
2. **Referenced Targets** (medium priority): Targets mentioned in other wikilinks
3. **Directory Completions** (lowest priority): Path-based suggestions

### File Creation

When navigating to a non-existent target:

1. **Path Resolution**: Determine appropriate file location
2. **Directory Creation**: Create parent directories if needed
3. **File Generation**: Create basic markdown file with title
4. **Navigation**: Open the newly created file

### Ambiguity Resolution

For targets matching multiple files:

1. **Detection**: Index identifies conflicting matches
2. **Diagnostics**: Error highlighting shows the issue
3. **Code Actions**: Quick fixes offer specific file selections
4. **Resolution**: User selects intended target with display text

## Integration Points

### Workspace Manager
- File discovery and indexing
- Change notifications for real-time updates
- Path resolution for target matching

### Parser System
- AST-based wikilink extraction
- Position tracking for LSP features
- Support for display text syntax

### LSP Protocol
- Completion provider registration
- Definition provider capabilities
- Diagnostic publishing
- Code action support

This system provides a comprehensive wikilink implementation that integrates seamlessly with standard LSP-capable editors.
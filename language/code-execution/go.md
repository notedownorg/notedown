# Go Code Execution

Go code blocks are executed by merging all `go` code blocks in a document into a single `main.go` file and running it with `go run`.

## Requirements

- The `go` binary must be available in the system PATH
- Go version 1.18 or later (for workspace mode support)

## Execution Model

1. **Collection**: All code blocks marked with `go` language identifier are collected in document order
2. **Merging**: Code blocks are combined into a single `main.go` file by:
   - Adding `package main` at the top
   - Concatenating all code blocks in document order
3. **Execution**: The merged file is written to a temporary directory and executed using `go run main.go` with the working directory set to the Notedown workspace root

## Code Block Format

Standard fenced code blocks with `go` language identifier:

```markdown
````go
fmt.Println("Hello, world!")
````
```

## Import Handling

- Users must structure code blocks so imports come before other declarations
- All code blocks are concatenated as-is (no automatic import extraction)
- Users are responsible for valid Go program structure

## Example

Given a Notedown document:

```markdown
# My Go Program

First, let's define a helper function:

```go
import (
    "fmt"
    "os"
    "strings"
)

func readFirstLine(filename string) (string, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return "", err
    }
    lines := strings.Split(string(content), "\n")
    if len(lines) > 0 {
        return lines[0], nil
    }
    return "", nil
}
```

Now let's use it in our main function:

```go
func main() {
    line, err := readFirstLine("README.md")
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        return
    }
    fmt.Printf("First line: %s\n", line)
}
```
```

This would be merged into:

```go
package main

import (
    "fmt"
    "os"
    "strings"
)

func readFirstLine(filename string) (string, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return "", err
    }
    lines := strings.Split(string(content), "\n")
    if len(lines) > 0 {
        return lines[0], nil
    }
    return "", nil
}

func main() {
    line, err := readFirstLine("README.md")
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        return
    }
    fmt.Printf("First line: %s\n", line)
}
```

## Execution Environment

- **Working Directory**: Set to the Notedown workspace root
- **File Access**: Full file system access (no restrictions)
- **Network Access**: Available (subject to system firewall/security settings)
- **Timeout**: Configurable timeout to prevent infinite loops (default: 30 seconds)


## Limitations

### Current Implementation
- **Standard Library Only**: Currently only Go standard library packages are supported
- **No External Dependencies**: No support for third-party modules or `go.mod` dependencies
- **Single File Programs**: Only supports single-file Go programs (package main only)
- **No Complex Projects**: No support for multiple files or complex project structures
- **No Interactive Input**: No support for interactive input during execution
- **Valid Go Structure**: Code blocks must form a valid Go program when merged
- **Explicit Functions**: Users must explicitly write `func main()` and any other required functions

### Future Extensions
Future versions may include:
- **Module Support**: Support for `go.mod` files and external dependencies
- **Multi-file Projects**: Support for complex project structures with multiple Go files
- **Package Development**: Support for creating and testing Go packages beyond package main
- **Workspace Dependencies**: Ability to import local packages from the workspace
- **Build Configuration**: Custom build flags and compilation options

# Code Block Execution

Notedown supports execution of code blocks within Markdown documents. This feature allows users to write executable code inline with their documentation and see the results immediately.

## Overview

Code block execution works by:
1. Collecting all code blocks of a specific language from a document
2. Merging them into a single executable file  
3. Running the code with the workspace root as the working directory
4. Capturing and displaying the output at the bottom of the file

## General Execution Model

### Collection and Merging
- All code blocks marked with the same language identifier are collected in document order
- Code blocks are combined into a single executable file using language-specific templates
- The system handles language-specific requirements like imports, package declarations, etc.

### Execution Environment
- **Working Directory**: Set to the Notedown workspace root
- **File Access**: Full read/write access to the workspace directory
- **Timeout**: Configurable timeout to prevent infinite loops (default: 30 seconds)

### Security Considerations
- Code execution runs with the same permissions as the Notedown process
- No sandboxing is provided - users should be cautious with untrusted code
- Full file system access available (code runs from workspace directory)
- Network access is unrestricted

## Supported Languages

- **[Go](go.md)** - Merge code blocks into main.go and execute with `go run`

## Future Extensions

This specification focuses on establishing the general framework. Future versions may include:

- Additional language support (Python, JavaScript, Rust, etc.)
- Sandboxing and security restrictions
- Interactive input/output
- Code block dependencies and execution order

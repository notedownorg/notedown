# Code Blocks Test

This file contains various code blocks for testing folding functionality.

## JavaScript Example

```javascript
function greet(name) {
    console.log(`Hello, ${name}!`);
    return true;
}

const users = ['Alice', 'Bob', 'Charlie'];
users.forEach(user => {
    greet(user);
});
```

## Python Example

```python
def fibonacci(n):
    if n <= 1:
        return n
    else:
        return fibonacci(n-1) + fibonacci(n-2)

# Generate first 10 fibonacci numbers
for i in range(10):
    print(f"F({i}) = {fibonacci(i)}")
```

## Go Example

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Starting timer...")
    
    for i := 5; i > 0; i-- {
        fmt.Printf("Countdown: %d\n", i)
        time.Sleep(1 * time.Second)
    }
    
    fmt.Println("Blast off! 🚀")
}
```

## Mixed Content

Here's some text between code blocks.

```bash
#!/bin/bash
echo "Setting up development environment..."

# Install dependencies
npm install

# Run tests
npm test

# Start development server
npm run dev
```

Final text after all code blocks.
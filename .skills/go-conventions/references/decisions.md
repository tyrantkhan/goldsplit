# Go Style Decisions

Source: https://google.github.io/styleguide/go/decisions

## Naming Conventions

### Underscores in Names
Names should generally avoid underscores, with three exceptions:
- Package names imported only by generated code
- Test, benchmark, and example function names in `*_test.go` files
- Low-level libraries interfacing with the OS or cgo

### Package Names
Packages must use lowercase letters and numbers only. Multi-word names remain unbroken and lowercase (e.g., `tabwriter` not `tab_writer`). Avoid generic names like `util`, `common`, or `helper`.

### Receiver Names
Method receivers should be:
- Short (one or two letters)
- Abbreviations of the type
- Applied consistently

```go
// Good
func (t Tray) Method() {}

// Bad
func (tray Tray) Method() {}
```

### Constant Names
Use MixedCaps for constants. Name based on role, not value:

```go
// Good
const MaxPacketSize = 512

// Bad
const MAX_PACKET_SIZE = 512
const kMaxBufferSize = 1024
```

### Initialisms
Maintain consistent casing within initialisms. For `URL`, use either `URL` or `url`, never `Url`.

### Getters
Omit `Get` prefixes; use the noun directly. For expensive operations, use `Compute` or `Fetch`:

```go
// Good: Counts() instead of GetCounts()
func (r *Record) Counts() int { ... }
```

### Variable Names
Name length should reflect scope size and usage frequency:
- Single-word names like `count` work well
- Add words to disambiguate (`userCount` vs `projectCount`)
- Omit type information (`users` not `userSlice`)
- Avoid repetition with surrounding context

### Repetitive Naming
Reduce redundancy between package and exported symbol names:

```go
// Good
package widget
func New() Widget { ... }

// Avoid
package widget
func NewWidget() Widget { ... }
```

## Error Handling

### Returning Errors
Functions that can fail should return `error` as the last parameter:

```go
func Process() error { ... }
func Lookup() (*Result, error) { ... }
```

Always return the `error` type, never concrete error types.

### Error Strings
Error messages should be lowercase and unpunctuated:

```go
// Good
err := fmt.Errorf("connection failed")

// Bad
err := fmt.Errorf("Connection Failed.")
```

### Handling Errors
Always handle errors deliberately. Avoid discarding with `_` unless documented why it's safe.

### In-Band Errors
Avoid returning special values like `-1` or empty strings to signal errors:

```go
// Good
func Lookup(key string) (string, bool) { ... }

// Bad
func Lookup(key string) string { ... }
```

### Indenting Error Flow
Check errors early and handle them before proceeding:

```go
// Good
if err != nil {
    return err
}
// normal code

// Bad
if err != nil {
    // handle
} else {
    // normal code
}
```

## Formatting and Structure

### Literal Formatting
Use composite literal syntax with field names:

```go
// Good
r := csv.Reader{
    Comma:   ',',
    Comment: '#',
}

// Bad
r := csv.Reader{',', '#', 4, false, false, false, false}
```

### Nil Slices
Prefer `nil` initialization for empty slices:

```go
// Good
var t []string

// Bad
t := []string{}
```

### Function Formatting
Keep function signatures on a single line.

### Conditionals
Don't break `if` statements across lines; extract variables:

```go
// Good
inTransaction := db.CheckStatus()
keysMatch := db.Compare(key1, key2)
if inTransaction && keysMatch {
    // ...
}
```

## Package Organization

### Imports
Organize into groups:
1. Standard library
2. Third-party and project packages
3. Protocol buffer imports
4. Side-effect imports

```go
import (
    "fmt"
    "os"

    "github.com/external/package"

    _ "path/to/sideeffect"
)
```

### Blank Imports
Only use `import _` in main packages or tests.

### Dot Imports
Never use `import .` syntax.

### Import Renaming
Rename imports only when necessary to avoid collisions.

## Commentary

### Doc Comments
All exported top-level names require doc comments starting with the name:

```go
// A Request represents a request to run a command.
type Request struct { ... }

// Encode writes JSON encoding of req to w.
func Encode(w io.Writer, req *Request) { ... }
```

### Package Comments
Place package comments immediately above the package clause:

```go
// Package math provides basic constants and mathematical functions.
package math
```

## Language Features

### Interfaces
Define interfaces in the package that *consumes* them. Return concrete types, not interfaces.

### Generics
Use generics judiciously. Start with concrete types.

### Receiver Types
- **Must use pointer:** Mutations needed, contains `sync.Mutex`, contains pointers to mutable data
- **Value is better:** Slices, small plain data, built-in types

### Synchronous Functions
Prefer synchronous functions. Callers can add concurrency.

### Goroutine Lifetimes
Make goroutine exit conditions explicit. Use context cancellation and `sync.WaitGroup`.

### Panic and Must Functions
Don't use `panic` for normal error handling. `MustXYZ` naming signals functions that panic on failure—only use during initialization.

## Testing

### Useful Test Failures
Format: `YourFunc(%v) = %v, want %v`

### Avoiding Assertion Libraries
Use `cmp.Equal` and `cmp.Diff`:

```go
if !cmp.Equal(got, want) {
    t.Errorf("BlogPost = %v, want = %v", got, want)
}
```

### Table-Driven Tests

```go
tests := []struct {
    input      string
    wantPieces []string
    wantErr    error
}{
    {
        input:      "1.2.3.4",
        wantPieces: []string{"1", "2", "3", "4"},
    },
    {
        input:   "hostname",
        wantErr: ErrBadHostname,
    },
}
```

## Common Libraries

### Contexts
`context.Context` should be the first parameter. Never add context as a struct field.

```go
func F(ctx context.Context, other string) {}
```

### crypto/rand
Always use `crypto/rand` for keys and security-sensitive randomness.

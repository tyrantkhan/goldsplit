# Google Go Style Guide - Complete Reference

Source: https://google.github.io/styleguide/go/

## Style Principles

### Clarity
The code's purpose and rationale is clear to the reader. Prioritize reader comprehension over author convenience.

Strategies:
- Use descriptive names
- Add comments explaining "why", not "what"
- Use whitespace effectively
- Refactor into modular functions

### Simplicity
The code accomplishes its goal in the simplest way possible.

- Code should be straightforward from top to bottom
- Avoid unnecessary abstractions and indirection
- **Least Mechanism principle**: Use the simplest available tool—prefer language constructs, then standard library, then core libraries before adding external dependencies

### Concision
The code has a high signal-to-noise ratio.

- Eliminate repetitive code patterns
- Use idiomatic Go constructs
- Choose clear, concise names
- Remove extraneous syntax

### Maintainability
The code is written such that it can be easily maintained.

- Easy for future programmers to modify safely
- APIs that scale gracefully
- Assumptions made explicit
- Minimize coupling and unused features
- Comprehensive tests with clear diagnostics

### Consistency
The code is consistent with the broader codebase.

Consistency within a package takes priority, but should not override documented style principles or global consistency.

---

## Naming Conventions

### MixedCaps
Use `MixedCaps` or `mixedCaps` (camel case) instead of underscores:
- Exported: `MaxLength`
- Unexported: `maxLength`

### Package Names
- Lowercase letters and numbers only
- Multi-word names remain unbroken: `tabwriter` not `tab_writer`
- Avoid generic names: `util`, `common`, `helper`

### Receiver Names
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
Use MixedCaps, name based on role not value:

```go
// Good
const MaxPacketSize = 512

// Bad
const MAX_PACKET_SIZE = 512
const kMaxBufferSize = 1024
```

### Initialisms
Maintain consistent casing. For `URL`, use either `URL` or `url`, never `Url`.

### Getters
Omit `Get` prefixes; use the noun directly:

```go
// Good
func (r *Record) Counts() int { ... }

// Bad
func (r *Record) GetCounts() int { ... }
```

For expensive operations, use `Compute` or `Fetch`.

### Variable Names
Name length should reflect scope size:
- Small scopes: single letters suffice
- Larger scopes: descriptive names needed

Rules:
- Add words to disambiguate (`userCount` vs `projectCount`)
- Omit type information (`users` not `userSlice`)
- Avoid repetition with surrounding context

### Reduce Redundancy

```go
// Good
package widget
func New() Widget { ... }

// Bad
package widget
func NewWidget() Widget { ... }
```

---

## Error Handling

### Returning Errors
Functions that can fail should return `error` as the last parameter:

```go
func Process() error { ... }
func Lookup() (*Result, error) { ... }
```

Always return the `error` type, never concrete error types.

### Error Strings
Lowercase and unpunctuated (unless starting with proper nouns):

```go
// Good
err := fmt.Errorf("connection failed")

// Bad
err := fmt.Errorf("Connection Failed.")
```

### Handling Errors
Always handle errors deliberately:
- Address them
- Return them
- Call `log.Fatal` in exceptional cases

Avoid discarding with `_` unless documented why it's safe.

### Indenting Error Flow
Check errors early, avoid `else` blocks:

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
    // normal code (awkwardly indented)
}
```

### Error Wrapping

```go
// Use %w when callers need errors.Is()/errors.As()
return fmt.Errorf("failed to process: %w", err)

// Use %v at system boundaries or when hiding internals
return fmt.Errorf("operation failed: %v", err)
```

Place `%w` at the end of error strings.

---

## Formatting and Structure

### gofmt
All Go source files must follow `gofmt` output formatting.

### Line Length
No fixed maximum. Prefer refactoring over splitting lines.

### Nil Slices
Prefer `nil` initialization:

```go
// Good
var t []string

// Bad
t := []string{}
```

### Composite Literals
Specify field names for types from other packages:

```go
// Good
r := csv.Reader{
    Comma:   ',',
    Comment: '#',
}

// Bad
r := csv.Reader{',', '#', 4, false, false, false, false}
```

### Function Signatures
Keep on a single line to prevent indentation confusion.

### Conditionals
Don't break `if` statements across lines; extract variables:

```go
// Good
inTransaction := db.CheckStatus()
keysMatch := db.Compare(key1, key2)
if inTransaction && keysMatch {
    // ...
}

// Bad
if db.CheckStatus() &&
    db.Compare(key1, key2) {
    // ...
}
```

---

## Package Organization

### Imports
Organize into groups with blank lines:

```go
import (
    "fmt"
    "os"

    "github.com/external/package"
    "go.traackr.com/internal/pkg"
)
```

### Import Rules
- Never use `import .` (dot imports)
- Rename imports only to avoid collisions
- Side-effect imports (`import _`) only in main packages or tests

---

## Documentation

### Doc Comments
All exported top-level names require doc comments starting with the name:

```go
// A Request represents a request to run a command.
type Request struct { ... }

// Encode writes JSON encoding of req to w.
func Encode(w io.Writer, req *Request) { ... }
```

### Package Comments
Place immediately above the package clause:

```go
// Package math provides basic constants and mathematical functions.
package math
```

### Comment Line Length
Aim for ~80 characters on narrow terminals.

---

## Interfaces

- Define interfaces where they're **consumed**, not produced
- Return concrete types, not interfaces
- Keep interfaces small and focused
- Don't export test doubles

---

## Concurrency

### Context
- `context.Context` should be the first parameter
- Never add context as a struct field
- Use `context.Background()` only in entrypoints

```go
func F(ctx context.Context, other string) {}
```

### Goroutines
- Make exit conditions explicit
- Use context cancellation and `sync.WaitGroup`
- Prefer synchronous functions; let callers add concurrency

```go
func (w *Worker) Run(ctx context.Context) error {
    var wg sync.WaitGroup
    for item := range w.q {
        wg.Add(1)
        go func() {
            defer wg.Done()
            process(ctx, item)
        }()
    }
    wg.Wait()
    return nil
}
```

### Copying
Don't copy structs containing `sync.Mutex` or types with pointer methods.

---

## Receiver Types

Choose based on:

**Must use pointer:**
- Mutations needed
- Contains `sync.Mutex`
- Contains pointers to mutable data

**Value is better:**
- Slices (unless reslicing)
- Small plain data
- Built-in types

**Guideline:** Keep all methods for a type either all pointer or all value receivers.

---

## Panic and Must Functions

- Don't use `panic` for normal error handling
- Reserve for impossible conditions
- `MustXYZ` naming signals functions that panic on failure
- Only use during initialization

```go
var DefaultVersion = MustParse("1.2.3")

func MustParse(version string) *Version {
    v, err := Parse(version)
    if err != nil {
        panic(fmt.Sprintf("MustParse(%q) = _, %v", version, err))
    }
    return v
}
```

---

## Testing

### Test Failure Messages
Include:
- Function name
- Input values (if short)
- Actual result before expected

Format: `YourFunc(%v) = %v, want %v`

### Comparison
Use `cmp.Equal` and `cmp.Diff`:

```go
if diff := cmp.Diff(got, want); diff != "" {
    t.Errorf("Process() mismatch (-got +want):\n%s", diff)
}
```

### Table-Driven Tests

```go
tests := []struct {
    name       string
    input      string
    wantPieces []string
    wantErr    error
}{
    {
        name:       "valid IP",
        input:      "1.2.3.4",
        wantPieces: []string{"1", "2", "3", "4"},
    },
    {
        name:    "invalid hostname",
        input:   "hostname",
        wantErr: ErrBadHostname,
    },
}

for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
        got, err := Parse(tc.input)
        // assertions...
    })
}
```

### t.Error vs t.Fatal
- Use `t.Error` to report multiple failures
- Use `t.Fatal` only when subsequent checks would be meaningless
- Never call `t.Fatal` from goroutines

---

## Common Libraries

### Flags
Flag names use underscores, variable names use MixedCaps:

```go
var pollInterval = flag.Duration("poll_interval", time.Minute, "...")
```

### crypto/rand
Always use `crypto/rand` for keys and security-sensitive randomness.

---

## Anti-Patterns to Avoid

1. Variable shadowing in nested scopes
2. Logging errors you return (let caller decide)
3. Panics in libraries (except API misuse detection)
4. Catching panics to avoid crashes
5. Ignoring errors without documentation
6. Using `math/rand` for security
7. Creating custom assertion libraries
8. Storing context in structs

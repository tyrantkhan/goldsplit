---
name: go-conventions
description: Go coding standards based on Google Go Style Guide. Use when writing, reviewing, or refactoring Go code. Apply for naming, error handling, package organization, and testing patterns.
metadata:
  author: traackr
  version: "1.0"
  source: https://google.github.io/styleguide/go/
---

# Go Conventions

Follow the Google Go Style Guide. For complete details:
- [Guide](references/guide.md) - Core style rules (canonical)
- [Decisions](references/decisions.md) - Detailed style choices
- [Best Practices](references/best-practices.md) - Common patterns

## Core Principles (in order)

1. **Clarity** - Reader comprehension over author convenience
2. **Simplicity** - Accomplish goals in the simplest way
3. **Concision** - High signal-to-noise ratio
4. **Maintainability** - Easy to modify safely
5. **Consistency** - Match the broader codebase

## Quick Reference

### Naming

- **MixedCaps** for all names, never underscores (`MaxLength`, `maxLength`)
- **Package names**: lowercase, single word, no underscores (`tabwriter` not `tab_writer`)
- **Avoid generic names**: no `util`, `common`, `helper`
- **Receivers**: short, 1-2 letters (`func (t Tray) Method()`)
- **No Get prefix**: use `Counts()` not `GetCounts()`
- **Initialisms**: consistent casing (`URL` or `url`, never `Url`)

### Error Handling

- Return `error` as last return value
- Error strings: lowercase, no punctuation (`"connection failed"`)
- Check errors early, avoid `else` blocks
- Use `%w` for wrapping when callers need `errors.Is()`/`errors.As()`
- Never discard errors with `_` unless documented why

### Formatting

- Run `gofmt` on all code
- No fixed line length, but refactor over splitting
- Prefer `nil` slices: `var t []string` not `t := []string{}`
- Keep function signatures on single lines
- Extract complex conditionals to named variables

### Imports

Organize in groups (blank line between):
1. Standard library
2. Third-party packages
3. Project packages

Never use `import .` (dot imports).

### Interfaces

- Define interfaces where they're **consumed**, not produced
- Return concrete types, accept interfaces
- Keep interfaces small and focused

### Testing

- Format: `YourFunc(%v) = %v, want %v`
- Use `cmp.Equal` and `cmp.Diff`, not custom assertions
- Table-driven tests with named struct fields
- Use `t.Error` to report multiple failures, `t.Fatal` only when test can't continue

### Concurrency

- Prefer synchronous functions; let callers add goroutines
- Make goroutine exit conditions explicit
- Use `context.Context` as first parameter
- Never store context in structs

### Things to Avoid

- `panic` for normal error handling (use `error` returns)
- Logging errors you also return (let caller decide)
- Variable shadowing in nested scopes
- Catching panics to avoid crashes
- `math/rand` for security (use `crypto/rand`)

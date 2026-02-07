# Go Best Practices

Source: https://google.github.io/styleguide/go/best-practices

## Function & API Design

### Naming Clarity
- Omit redundant information from function names
- Use verb-like names for actions, noun-like names for values
- For many parameters, use option structs or variadic options

### Option Structures vs. Variadic Options
- **Option structs**: When most callers need options, shared across functions
- **Variadic options**: When most callers need no configuration

## Error Handling

### Structured Errors Over String Matching
Callers should inspect errors using `errors.Is()`/`errors.As()` rather than parsing error messages.

### Error Wrapping Guidelines
- Use `%w` for preserving error chains when callers need programmatic inspection
- Use `%v` when transforming errors at system boundaries
- Place `%w` at the end of error strings

```go
// Preserves chain
return fmt.Errorf("failed to process: %w", err)

// Hides internals
return fmt.Errorf("operation failed: %v", err)
```

### Avoid Redundant Context
Don't add annotation if its sole purpose is to indicate failure without adding new information.

## Documentation

### Parameter Documentation
Only document non-obvious or error-prone parameters.

### Context Behavior
Documenting that context cancellation interrupts a function is unnecessary—it's implicit.

### Concurrency Safety
Assume read-only operations are safe for concurrent use. Explicitly document mutations that aren't thread-safe.

## Testing Patterns

### Avoid Assertion Helpers
Keep validation logic in the test function itself. If validation is needed across tests, return an `error` value.

### Table-Driven Tests with Field Names
Use struct field names, especially for cases spanning 20+ lines.

### Setup Scope
Call setup functions explicitly within tests that need them rather than centralizing in `TestMain`.

### Fatal vs. Error
- Use `t.Fatal` for setup failures or when the test cannot proceed
- Use `t.Error` in table entries when one case fails but others can continue
- Never call `t.Fatal` from goroutines

## Package Organization

### Avoid Overly Generic Names
Packages named `util`, `helper`, or `common` obscure what they provide.

### Test Helper Packages
Name test double packages by appending `test` to the production package name (e.g., `creditcardtest`).

## Anti-Patterns to Avoid

1. **Variable shadowing** in new scopes that creates confusion
2. **Logging errors you return** — let the caller decide
3. **Panics in libraries** — except for detecting API misuse
4. **Catching panics** to avoid crashes — propagates corrupted state
5. **Ignoring errors** except when orchestrating operations where only first error matters

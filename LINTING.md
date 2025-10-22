# Golangci-lint Configuration Guide

## Overview
Konfigurasi golangci-lint yang telah dioptimalkan untuk project `ewallet-ums` dengan mengikuti best practices Go development.

## Enabled Linters

### Bug Detection
- **bodyclose**: Checks whether HTTP response bodies are closed
- **errcheck**: Checks for unchecked errors
- **gosec**: Inspects source code for security problems
- **govet**: Reports suspicious constructs
- **staticcheck**: Advanced static analysis
- **typecheck**: Type checking

### Code Quality
- **deadcode**: Finds unused code (deprecated, replaced by unused)
- **ineffassign**: Detects ineffectual assignments
- **unparam**: Reports unused function parameters
- **unused**: Checks for unused constants, variables, functions and types

### Performance
- **gocritic**: Comprehensive Go source code linter
- **prealloc**: Finds slice declarations that could potentially be preallocated

### Style & Formatting
- **gofmt**: Checks whether code was gofmt-ed
- **goimports**: Checks import sorting and formatting
- **gofumpt**: Stricter gofmt
- **lll**: Reports long lines
- **misspell**: Finds commonly misspelled English words
- **whitespace**: Detects leading and trailing whitespace

### Complexity
- **gocognit**: Computes cognitive complexity
- **gocyclo**: Computes cyclomatic complexity
- **funlen**: Detects long functions
- **nestif**: Reports deeply nested if statements

### Error Handling
- **errname**: Checks that sentinel errors are prefixed with Err
- **errorlint**: Finds code that will cause problems with error wrapping

### Best Practices
- **godot**: Checks if comments end in a period
- **nilerr**: Finds code that returns nil even if it checks that error is not nil
- **nilnil**: Checks that there is no simultaneous return of nil error and invalid value
- **noctx**: Finds sending HTTP request without context.Context
- **contextcheck**: Checks whether functions use context.Context

## Configuration Highlights

### Line Length
```yaml
lll:
  line-length: 140
  tab-width: 1
```

### Magic Numbers
Magic numbers yang diizinkan:
- HTTP status codes: 100, 200, 400, 401, 403, 404, 500
- Common numbers: 0, 1, 2

### Function Complexity
```yaml
funlen:
  lines: 100
  statements: 50

gocyclo:
  min-complexity: 15
```

### Import Grouping
```yaml
goimports:
  local-prefixes: github.com/ibnuzaman/ewallet-ums
```

## Exclusions

### Test Files
Linters berikut dinonaktifkan untuk file `_test.go`:
- mnd (magic number detector)
- funlen (function length)
- goconst (repeated strings)
- gocyclo (cyclomatic complexity)
- errcheck (unchecked errors)
- dupl (code duplication)
- gosec (security checks)

### Special Files
- **main.go**: `gochecknoinits` disabled

## Usage

### Run Linter
```bash
make lint
# or
golangci-lint run
```

### Auto-fix Issues
```bash
make lint-fix
# or
golangci-lint run --fix
```

### Run All Checks
```bash
make check
```

### CI Pipeline
```bash
make ci
```

## Common Issues & Solutions

### Issue: Import grouping error
**Solution**: Use `goimports` to format imports
```bash
goimports -w .
```

### Issue: Struct field alignment
**Solution**: Reorder struct fields to minimize memory padding
```go
// Before (48 bytes)
type Response struct {
    Success bool        // 1 byte + 7 padding
    Message string      // 16 bytes
    Data    interface{} // 16 bytes
    RequestID string    // 16 bytes
}

// After (40 bytes)
type Response struct {
    Data      interface{} // 16 bytes
    RequestID string      // 16 bytes
    Message   string      // 16 bytes
    Success   bool        // 1 byte
}
```

### Issue: Comment must end in period
**Solution**: Add period at the end of comments
```go
// Bad
// LoggerMiddleware logs HTTP requests

// Good
// LoggerMiddleware logs HTTP requests.
```

## Integration with VS Code

Add to `.vscode/settings.json`:
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--fast"
  ],
  "go.lintOnSave": "workspace",
  "editor.formatOnSave": true,
  "go.formatTool": "goimports"
}
```

## CI/CD Integration

### GitHub Actions
```yaml
- name: golangci-lint
  uses: golangci/golangci-lint-action@v3
  with:
    version: latest
    args: --timeout=5m
```

## Best Practices

1. **Run linter before commit**
   ```bash
   make check
   ```

2. **Fix issues incrementally**
   ```bash
   golangci-lint run --new
   ```

3. **Use nolint sparingly**
   ```go
   // Only when absolutely necessary
   //nolint:errcheck // Reason: explanation why it's safe to ignore
   _ = file.Close()
   ```

4. **Keep configuration updated**
   - Review linter updates quarterly
   - Remove deprecated linters
   - Add new relevant linters

## Performance Tips

- Use `--fast` flag for quick checks during development
- Enable `--allow-parallel-runners` in CI
- Set appropriate timeout (5m default)
- Use `--new` to check only new code

## Resources

- [golangci-lint Documentation](https://golangci-lint.run/)
- [Enabled Linters](https://golangci-lint.run/usage/linters/)
- [Configuration Reference](https://golangci-lint.run/usage/configuration/)

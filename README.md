# spvalidator

`spvalidator` is a small validation library for Go with direct function-call validators.

```go
package main

import "github.com/yorukot/spvalidator"

func main() {
	if err := spvalidator.Email("user@example.com"); err != nil {
		panic(err)
	}
}
```

Validators return `nil` on success and `*spvalidator.ValidationError` on failure.

```go
if err := spvalidator.Required(""); spvalidator.IsValidationError(err) {
	// handle validation error
}
```

Field validators use exported struct field paths.

```go
type Signup struct {
	Password string
	Confirm  string
}

_ = spvalidator.EqField(Signup{"secret", "secret"}, "Password", "Confirm")
```

The package is independently implemented and aims for practical validation coverage for the supported tags. Very broad standards such as postcode formats, BCP 47, HTML detection, cron syntax, and RFC-perfect URL/domain parsing use documented practical checks rather than exhaustive standard conformance.

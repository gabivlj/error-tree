# Rrors

A proof-of-concept Golang error library that renders your error stack trace as a tree. Similar to xerrors but with support to `errors.Join`.

See the tests folder to see sample renderings.

# Is this ready?

This is a proof of concept and hacky code all around, I'm just testing the idea :)

# Usage

This is at your own risk, code is not perfect, but the idea is wrapping errors and preserving the stacktrace like `xerrors` did, but also
adaption to errors.Join from go 1.20

```golang
err := rrors.Errorf("some error: %w", errors.Join(...))
// prints tree with stacktrace
fmt.Println(err)
```

## Untested things

- rrors.Errorf("%w %w") ?
- errors.Is(rrors.Errorf(...)) ?
- errors.As(...)

## Plans

Add ways to query an error tree like doing:

- `rerrs.CountErrors(err, os.ErrNotFound)` // counts the number of ErrNotFound in the error tree
- `rerrs.AsMultiple[MyErrorType](err)` // returns errors of this type in the error tree

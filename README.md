# quotecheck

[![tag](https://img.shields.io/github/tag/peczenyj/fmtquotecheck.svg)](https://github.com/peczenyj/fmtquotecheck/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-%23007d9c)
[![GoDoc](https://pkg.go.dev/badge/github.com/peczenyj/fmtquotecheck)](http://pkg.go.dev/github.com/peczenyj/fmtquotecheck)
[![Go](https://github.com/peczenyj/fmtquotecheck/actions/workflows/go.yml/badge.svg)](https://github.com/peczenyj/fmtquotecheck/actions/workflows/go.yml)
[![Lint](https://github.com/peczenyj/fmtquotecheck/actions/workflows/lint.yml/badge.svg)](https://github.com/peczenyj/fmtquotecheck/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/peczenyj/fmtquotecheck/graph/badge.svg?token=9y6f3vGgpr)](https://codecov.io/gh/peczenyj/fmtquotecheck)
[![Report card](https://goreportcard.com/badge/github.com/peczenyj/fmtquotecheck)](https://goreportcard.com/report/github.com/peczenyj/fmtquotecheck)
[![CodeQL](https://github.com/peczenyj/fmtquotecheck/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/peczenyj/fmtquotecheck/actions/workflows/github-code-scanning/codeql)
[![Dependency Review](https://github.com/peczenyj/fmtquotecheck/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/peczenyj/fmtquotecheck/actions/workflows/dependency-review.yml)
[![License](https://img.shields.io/github/license/peczenyj/fmtquotecheck)](./LICENSE)
[![Latest release](https://img.shields.io/github/release/peczenyj/fmtquotecheck.svg)](https://github.com/peczenyj/fmtquotecheck/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/peczenyj/fmtquotecheck.svg)](https://github.com/peczenyj/fmtquotecheck/releases/latest)
[![Last commit](https://img.shields.io/github/last-commit/peczenyj/fmtquotecheck.svg)](https://github.com/peczenyj/fmtquotecheck/commit/HEAD)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/peczenyj/fmtquotecheck/blob/main/CONTRIBUTING.md#pull-request-process)

Verify when safely escape and single quote strings on fmt.Sprintf.

This code is based on [perfsprint](https://github.com/catenacyber/perfsprint) analyzer.

## Motivation

While it seems safe write `'%s'` to just delimit some string, it may have nasty consequences, like if the string already contains a quote char it add some surprises in our output.
But go supports a better alternative: `%q` is [defined](https://pkg.go.dev/fmt) as `a single-quoted character literal safely escaped with Go syntax.`

## Instruction

```sh
go install github.com/peczenyj/fmtquotecheck/cmd/fmtquotecheck@latest
```

## Usage

```go
package main

import "fmt"

func main(){
    fmt.Printf("hello '%s'", "world") // we should use %q here 
}
```

```console
$ fmtquotecheck ./main.go 
./main.go:6:16: explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf
```

by using the option `-fix` the linter will convert all `'%s'` to `%q`.

## CI

### CircleCI

```yaml
- run:
    name: install fmtquotecheck
    command: go install github.com/peczenyj/fmtquotecheck/cmd/fmtquotecheck@latest

- run:
    name: run fmtquotecheck
    command: fmtquotecheck ./...
```

### GitHub Actions

```yaml
- name: install fmtquotecheck
  run: go install github.com/peczenyj/fmtquotecheck/cmd/fmtquotecheck@latest

- name: run fmtquotecheck
  run: fmtquotecheck ./...
```

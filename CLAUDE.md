# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

**dish** is a dynamically-typed programming language (interpreter written in Go) for processing input data — especially JSON — with traditional operators instead of jq-style queries. Every data type is an object, all symbolic operators are shorthand for object methods, and there are no control-flow keywords (`if`/`for`/`while` are replaced by methods like `.switch`, `.times`, `.while`).

The whole interpreter is a flat `package main` at the repo root, with no dependencies beyond the Go standard library.

## Build, run, test

```sh
go build -o /usr/local/bin/dish       # or: go build -o dish .

dish path/to/file.dish                # run a file
dish -e '"Hello World!"'              # run an expression
dish -p -e 'stdin.map:d(d.name)'      # -p pretty (colorized), -f format (indented JSON)
dish -d -e '1 + 2'                    # -d dumps the parsed token stream to STDERR and exits
echo '[1,2,3]' | dish -e 'stdin.sum'  # STDIN is JSON-parsed into `stdin`/`$_`, else split into an Array of Strings
```

Flags may be bundled (`-pe`, `-fe`); the `e` must be last when bundled. Non-flag trailing arguments land in `argv` / `$0` / `$1`…

```sh
go test ./...                # full suite — must be run from the repo root (test paths are relative)
go test -run TestDot ./...   # one test file
go test -v ./...             # shows the "tests/x.dish passed N of M tests" lines
```

`go vet ./...` is clean and worth keeping that way, since it is only useful as a gate while it exits 0. Run it separately from `go test`, which applies just a subset of vet's checks.

The `go 1.16` directive in `go.mod` reflects the oldest toolchain that can build this — the code uses no newer language or stdlib features.

If a sandbox blocks the default build cache, set `GOCACHE` to a writable directory.

## Test file format

Each `tests/*.dish` file is a single Array block whose entries are expressions, each annotated with a `##` comment holding the expected output:

```
[
    [1, 2, 3] | [4, 2, 1]   ## [1, 2, 3, 4]
    'binary' | :c(c.vowel)  ## true
]
```

`test()` in `dish.go` runs the file and compares evaluated array item `i` against lexed comment `i` — **positionally**. Any `##` comment that is not an expectation (a header comment, a commented-out line) shifts every subsequent expectation out of alignment and produces bogus failures. Adding a new `tests/<name>.dish` also requires a matching `func Test<Name>` in `dish_test.go`; nothing discovers test files automatically.

`examples/` holds Rosetta Code–style programs. They are documentation/smoke-test material and are **not** executed by `go test`.

## Architecture

The pipeline is four stages, one file each, followed by per-operator implementation files:

1. **`chomp.go` — Lexer.** `Lexify()` produces `*Token`s tagged `NUM`/`STR`/`VAR`/`OP1` (unary prefix)/`OP2` (binary)/`OPX` (alphanumeric method name)/`MEM`/`BLK`/`FIN`/`EOF`. Whether an operator is unary or binary is decided here from surrounding context, so `-5` and `a - 5` lex differently. This file also owns `Token.Precedence()` (the operator precedence table), `Higher()`, and the block-matching helpers. `##` comments are collected into `Lexer.coms` for the test harness.
2. **`churn.go` — Parser.** `Parser.Churn()` is a shunting-yard loop converting the token stream into RPN stored on `Block.toks`. Blocks (`(...)`, `[...]`, `{...}`) and short-circuit right-hand operands become nested `*Block`s so they can be evaluated lazily. Also holds `Block`, `Run`, assignment/scoping (`Assign`, `Variable.Assign`), and string interpolation.
3. **`chirp.go` — Types + evaluator.** Defines the six value types (`Hash`, `Array`, `String` = `[]rune`, `Number` = `*big.Rat` plus an `inf` flag, `Boolean`, `Null`) plus `Variable`. `Block.Chirp()` walks the RPN with a value stack (`Register`/`Deregister`) and dispatches every operator.
4. **`dish.go` — Entry point.** Flag parsing, STDIN binding (`bind`/`parse`), output printing (raw scalar vs. JSON for Array/Hash, optional ANSI colorization), and the `test()` helper used by `dish_test.go` (which is why non-test `dish.go` imports `testing`).

### Operator files and their naming convention

Each remaining `.go` file implements one operator symbol and is named after that symbol's nickname:

| File | Symbol | File | Symbol | File | Symbol |
| --- | --- | --- | --- | --- | --- |
| `cross.go` | `+` | `dash.go` | `-` | `splat.go` | `*` |
| `split.go` | `/` | `grep.go` | `%` | `pow.go` | `^` |
| `sharp.go` | `#` | `thump.go` | `@` | `dot.go` | `.` |
| `boom.go` | `&` | `wham.go` | `\|` | `twiddle.go` | `~` |
| `whiz.go` | `?` | `bang.go` | `!` | `bars.go` | `=` |
| `wiki.go` | `<` | `waka.go` | `>` | `criss.go` | `==` / coercions |
| `snap.go` | Number methods | `snip.go` | String methods | | |

Method names inside those files follow the same scheme:

- `X` — the binary form (`Splat` = `*`).
- `TopX` — the unary prefix form (`TopSplat` = prefix `*`). Note `TopHat` is prefix `^`.
- `DoubleX` — the doubled symbol (`DoubleSplat` = `**`, `DoubleSplit` = `//`, `DoubleGrep` = `%%`, `DoubleBoom` = `&&`, `DoubleWham` = `||`, `TwiddleDee` = `~~`, `DoubleDot` = `..`, `DoubleWhiz` = `??`).
- `WikiWiki`/`WakaWaka` = binary `<<`/`>>`; `WikiBars`/`WakaBars` = `<=`/`>=`.

Below those, per-type helpers spell out the actual behavior (`MapHash`, `IntersectArray`, `RotateString`, …) — this is where the semantics documented in `tests/README.md` live.

### Dispatch and aliases — the thing to get right

The `switch t.lit` blocks in `chirp.go` (`Chirp()`, `OP1` case ≈ lines 423–503, `OP2` case ≈ lines 556–712) are the single source of truth mapping a token literal to a handler. Both the symbol and every word alias for it are listed there (`case "*", "multiply", "mult", "repeat", "times", "map", "each", "join", "dot":`).

Each handler then re-checks `t.lit` against the aliases legal *for that operand type pairing* and calls `t.TypeMismatch(x, y)` otherwise:

```go
if t.lit != "*" && t.lit != "*=" && t.lit != "map" && t.lit != "each" {
    t.TypeMismatch(x, y)
}
```

So **adding or renaming an alias means editing two places**: the `case` list in `chirp.go` *and* every guard in the operator file that should accept it. Missing the second half yields a runtime type-mismatch panic instead of the new alias working. Compound assignment operators (`+=`, `*=`, …) route through the same handlers, which is why guards typically also allow the `=`-suffixed literal.

Two more invariants every handler observes:

- **Unwrap first.** `*Block` and `*Variable` flow through the evaluator as values, so every handler begins with `case *Block: return t.X(x.Run(), b)` and `case *Variable: return t.X(x.Value(), b)` before matching real types. New handlers must do the same or they will silently mishandle variables and lazily-evaluated blocks.
- **Booleans and Nulls coerce to Numbers** (`1`/`0`) in numeric context; other type pairings that are not explicitly handled fall through to `TypeMismatch`, which panics.

### Scoping and re-entrancy

`Block` carries `src` (lexical parent), `args`/`def` (parameter names and defaults), and a stack of `*Run` frames with `cur` pointing at the active one — that stack is what makes recursion and repeated block invocation work. `NewBlock()` seeds the global frame with `true`, `false`, `null`, `inf`, `stdin`/`$_`, `argv`/`$0`, and `$1`…`$n`. `:` defines a variable in local scope, `=` assigns in the enclosing/global scope.

## Documentation is part of the contract

`README.md` documents the language (syntax, blocks, member access, interpolation) and `tests/README.md` is the full operator reference — every symbol, its operand type pairings, its method-name aliases, and a worked example with the expected result. When adding an operator, alias, or behavior, update the relevant table in `tests/README.md` and add expectations to the matching `tests/*.dish` file; the two are meant to agree line for line.

## Style

Go sources use 4-space indentation (no tabs), `import("a"; "b")` single-line grouping for short import lists, and braces-with-spaces literals like `Number{ val: ... }` / `Array{ }`. Match the surrounding file rather than running `gofmt` over it — the codebase is deliberately not gofmt-formatted, and reformatting produces enormous unrelated diffs.

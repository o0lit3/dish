package main

import (
    "os"
    "sort"
    "strconv"
    "fmt"
    "bytes"
    "sync"
    "time"
    "regexp"
    "context"
    "os/exec"
    "strings"
    "testing"
    "runtime"
    "path/filepath"
)

// A combination sweep over operand types and syntax, run as one dish process per
// expression. It is gated behind SWEEP so that `go test ./...` stays fast:
//
//     SWEEP=1 go test -run TestSweep -v ./...
//
// Process isolation is required rather than incidental. An interpreter that loops
// forever or overflows the Go stack cannot be recovered in-process, so those two
// classes of defect are only observable as a killed child process.
//
// Only unambiguous defects fail the test: a Go-level crash, a hang, a spacing
// variant that changes behavior, or an expression matching neither of its two
// parenthesizations. A dish error carrying a position is a correct outcome. The
// remaining axes log what they find, because judging them needs a human.

var values = []string{ "0", "-2", "3", "inf", "(0-inf)", "null", "true", "\"\"", "\"ab\"", "[]", "[1,2]", "[9]", "{}", "{a:1}" }

var binaries = []string{ "+", "-", "*", "/", "%", "^", "#", "@", "&", "|", "~", "?", "!", "<", ">",
    "==", "!=", "<=", ">=", "<<", ">>", "**", "//", "%%", "&&", "||", "~~", "..", "??" }

var unaries = []string{ "+", "-", "*", "/", "%", "^", "#", "@", "&", "|", "~", "?", "!",
    "++", "--", "**", "//", "%%", "<<", ">>" }

var sample = []string{ "3", "(-2)", "0", "\"ab\"", "[1,2]", "{a:1}", "null", "true", "[[1,2],[3]]", "[{a:1},{b:2}]" }

var objects = []string{ "[1,2,3]", "\"abc\"", "12", "{a:1}", "null", "true", "[]" }

var members = []string{ "0", "1", "2", "-1", "a", "\"a\"", "sum", "len", "str", "(1)", "[0,1]" }

var placed = regexp.MustCompile(` at \d+:\d+`)

type group struct {
    unary bool
    sym string
    words []string
}

type answer struct {
    kind string
    text string
}

func compile(test *testing.T) string {
    bin := filepath.Join(test.TempDir(), "dish")

    if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
        test.Fatalf("cannot build dish: %v\n%s", err, out)
    }

    return bin
}

// evaluate classifies one expression. Positions are stripped from dish errors so
// that callers can compare behavior; a shifted column is a correct difference.
func evaluate(bin string, expr string) answer {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    var stdout, stderr bytes.Buffer

    cmd := exec.CommandContext(ctx, bin, "-e", expr)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    msg := strings.TrimSpace(stderr.String())

    if lines := strings.SplitN(msg, "\n", 2); len(lines) > 0 {
        msg = lines[0]
    }

    switch {
    case ctx.Err() == context.DeadlineExceeded:
        return answer{ kind: "hang" }
    case err == nil:
        return answer{ kind: "ok", text: strings.TrimSpace(stdout.String()) }
    case strings.Contains(msg, "internal error"), strings.Contains(msg, "goroutine"), strings.Contains(msg, "panic"):
        return answer{ kind: "crash", text: msg }
    case strings.HasPrefix(msg, "dish: "):
        return answer{ kind: "error", text: placed.ReplaceAllString(msg, "") }
    }

    return answer{ kind: "crash", text: msg }
}

// resolve evaluates every expression concurrently and returns them keyed by source.
func resolve(bin string, exprs []string) map[string]answer {
    uniq := []string{ }
    seen := map[string]bool{ }

    for _, expr := range exprs {
        if !seen[expr] {
            seen[expr] = true
            uniq = append(uniq, expr)
        }
    }

    out := map[string]answer{ }
    gate := make(chan bool, runtime.NumCPU())
    lock := sync.Mutex{ }
    wait := sync.WaitGroup{ }

    for _, expr := range uniq {
        wait.Add(1)

        go func(expr string) {
            defer wait.Done()
            gate <- true
            got := evaluate(bin, expr)
            <-gate

            lock.Lock()
            out[expr] = got
            lock.Unlock()
        }(expr)
    }

    wait.Wait()
    return out
}

func broken(test *testing.T, got map[string]answer) int {
    bad := 0

    for expr, val := range got {
        if val.kind == "crash" || val.kind == "hang" {
            test.Errorf("%s: %s %s", expr, val.kind, val.text)
            bad++
        }
    }

    return bad
}

func TestSweep(test *testing.T) {
    if os.Getenv("SWEEP") == "" {
        test.Skip("set SWEEP=1 to run the combination sweep")
    }

    bin := compile(test)

    // Every value pairing against every binary operator.
    test.Run("Types", func(test *testing.T) {
        var exprs []string

        for _, a := range values {
            for _, b := range values {
                for _, op := range binaries {
                    exprs = append(exprs, fmt.Sprintf("%s %s %s", a, op, b))
                }
            }
        }

        got := resolve(bin, exprs)
        errs := 0

        for _, val := range got {
            if val.kind == "error" {
                errs++
            }
        }

        fmt.Printf("sweep types: %d combinations, %d dish errors, %d problems\n",
            len(got), errs, broken(test, got))
    })

    // Every prefix operator against every value.
    test.Run("Unary", func(test *testing.T) {
        var exprs []string

        for _, op := range unaries {
            for _, a := range values {
                exprs = append(exprs, op + a)
            }
        }

        got := resolve(bin, exprs)

        fmt.Printf("sweep unary: %d expressions, %d problems\n", len(got), broken(test, got))
    })

    // Whitespace around "." is insignificant, so the three spaced forms of a member
    // access must agree. The unspaced form differs only when both sides are numeric,
    // where "12.1" is a floating point number rather than member access.
    test.Run("Spacing", func(test *testing.T) {
        var exprs []string
        forms := func(obj string, mem string) []string {
            return []string{ obj + "." + mem, obj + " ." + mem, obj + ". " + mem, obj + " . " + mem }
        }

        for _, obj := range objects {
            for _, mem := range members {
                exprs = append(exprs, forms(obj, mem)...)
            }
        }

        got := resolve(bin, exprs)
        split := 0

        for _, obj := range objects {
            for _, mem := range members {
                all := forms(obj, mem)

                for _, expr := range all[2:] {
                    if got[expr] != got[all[1]] {
                        test.Errorf("%q and %q disagree: %v vs %v", all[1], expr, got[all[1]], got[expr])
                        split++
                    }
                }
            }
        }

        fmt.Printf("sweep spacing: %d expressions, %d disagreements, %d problems\n",
            len(got), split, broken(test, got))
    })

    // "a OP1 b OP2 c" must evaluate as one of its two parenthesizations. Matching
    // neither means the operator table and the parser disagree.
    test.Run("Order", func(test *testing.T) {
        var exprs []string
        forms := func(one string, two string) []string {
            return []string{
                fmt.Sprintf("6 %s 3 %s 2", one, two),
                fmt.Sprintf("(6 %s 3) %s 2", one, two),
                fmt.Sprintf("6 %s (3 %s 2)", one, two),
            }
        }

        for _, one := range binaries {
            for _, two := range binaries {
                exprs = append(exprs, forms(one, two)...)
            }
        }

        got := resolve(bin, exprs)
        odd := 0

        for _, one := range binaries {
            for _, two := range binaries {
                all := forms(one, two)

                if got[all[0]] != got[all[1]] && got[all[0]] != got[all[2]] {
                    test.Errorf("%q matches neither grouping: %v, left %v, right %v",
                        all[0], got[all[0]], got[all[1]], got[all[2]])
                    odd++
                }
            }
        }

        fmt.Printf("sweep order: %d operator pairs, %d matching neither grouping, %d problems\n",
            len(binaries) * len(binaries), odd, broken(test, got))
    })

    // Word aliases against their symbol. Aliases in a chirp.go case group share a
    // handler, not a meaning -- the handler branches on the literal -- so a value
    // difference is not a defect on its own. Crashes are, and the rejections are
    // worth reading: that list is where a guard missing an alias shows up. It prints
    // sorted by name rather than by count so that two runs can be diffed -- a newly
    // added alias whose guard was forgotten arrives as a new line, and those have a
    // low case count precisely because they affect one operand pairing.
    test.Run("Aliases", func(test *testing.T) {
        groups := aliases(test)
        var exprs []string
        pairs := map[string]string{ }
        named := map[string]string{ }

        for _, grp := range groups {
            for _, word := range grp.words {
                for _, a := range sample {
                    if grp.sym == "" {
                        exprs = append(exprs, fmt.Sprintf("%s.%s", a, word))
                        continue
                    }

                    if grp.unary {
                        base := grp.sym + a
                        alias := fmt.Sprintf("%s.%s", a, word)
                        pairs[alias] = base
                        named[alias] = fmt.Sprintf(".%s (for %s)", word, grp.sym)
                        exprs = append(exprs, base, alias)
                        continue
                    }

                    for _, b := range sample {
                        base := fmt.Sprintf("%s %s %s", a, grp.sym, b)
                        alias := fmt.Sprintf("%s.%s(%s)", a, word, b)
                        pairs[alias] = base
                        named[alias] = fmt.Sprintf(".%s (for %s)", word, grp.sym)
                        exprs = append(exprs, base, alias)
                    }
                }
            }
        }

        got := resolve(bin, exprs)
        count := map[string]int{ }
        first := map[string]string{ }

        for alias, base := range pairs {
            if got[base].kind == "ok" && got[alias].kind == "error" {
                count[named[alias]]++

                if _, ok := first[named[alias]]; !ok {
                    first[named[alias]] = base
                }
            }
        }

        keys := []string{ }

        for key := range count {
            keys = append(keys, key)
        }

        sort.Strings(keys)

        fmt.Printf("sweep aliases: %d groups, %d checks, %d rejections where the symbol worked, %d problems\n",
            len(groups), len(pairs), len(keys), broken(test, got))

        for _, key := range keys {
            fmt.Printf("    %-24s %3d cases   e.g. %s\n", key, count[key], first[key])
        }
    })

    // Every rational dish can represent, it should compute exactly. A value that
    // passed through a 64-bit float betrays itself with a large power-of-two
    // denominator, so each operation that must be exact is asked for its ratio.
    test.Run("Exact", func(test *testing.T) {
        vals := []string{ "3", "-5", "(1/3)", "(2/7)", "0.1", "(3/4)", "7", "(2/3)" }
        ops := []string{ "+", "-", "*", "/", "//", "%", "%%" }
        var exprs []string

        for _, a := range vals {
            for _, b := range vals {
                for _, op := range ops {
                    exprs = append(exprs, fmt.Sprintf("%%(%s %s %s)", a, op, b))
                }
            }

            for _, pow := range []string{ "2", "3", "-2" } {
                exprs = append(exprs, fmt.Sprintf("%%(%s ^ %s)", a, pow))
            }

            exprs = append(exprs, fmt.Sprintf("%%(-%s)", a), fmt.Sprintf("%%(%s @ 3)", a),
                fmt.Sprintf("%%([%s, 2].sum)", a), fmt.Sprintf("%%([%s, 2].avg)", a),
                fmt.Sprintf("%%([%s, 2].product)", a))
        }

        got := resolve(bin, exprs)
        denom := regexp.MustCompile(`"denom": (\d+)`)
        loose := 0

        for expr, val := range got {
            if val.kind != "ok" {
                continue
            }

            found := denom.FindStringSubmatch(val.text)

            if found == nil {
                continue
            }

            den, err := strconv.ParseInt(found[1], 10, 64)

            if err != nil || den <= 1 << 20 || den & (den - 1) != 0 {
                continue
            }

            test.Errorf("%s lost exactness: denominator %d is a power of two", expr, den)
            loose++
        }

        fmt.Printf("sweep exact: %d expressions, %d inexact, %d problems\n",
            len(got), loose, broken(test, got))
    })
}

// aliases reads both dispatch tables out of chirp.go, pairing each group's symbol
// with the word aliases that share its handler. Groups are kept as a list because
// more than one group can lead with the same symbol.
func aliases(test *testing.T) []group {
    src, err := os.ReadFile("chirp.go")

    if err != nil {
        test.Fatalf("cannot read chirp.go: %v", err)
    }

    body := string(src)
    one := strings.Index(body, "case OP1:")
    two := strings.Index(body, "case OP2:")

    if one < 0 || two < 0 || two < one {
        test.Fatal("cannot find the OP1 and OP2 dispatch in chirp.go")
    }

    out := []group{ }
    cases := regexp.MustCompile(`(?m)^\s+case ((?:"[^"]*"(?:, )?)+):`)
    quoted := regexp.MustCompile(`"([^"]*)"`)
    worded := regexp.MustCompile(`^[a-z][a-z0-9]*$`)

    for _, part := range []struct{ unary bool; body string }{
        { unary: true, body: body[one:two] },
        { unary: false, body: body[two:] },
    } {
        for _, found := range cases.FindAllStringSubmatch(part.body, -1) {
            sym := ""
            var words []string

            for _, lit := range quoted.FindAllStringSubmatch(found[1], -1) {
                switch {
                case worded.MatchString(lit[1]):
                    words = append(words, lit[1])
                case sym == "" && !strings.HasSuffix(lit[1], "="):
                    sym = lit[1]
                }
            }

            if len(words) > 0 {
                out = append(out, group{ unary: part.unary, sym: sym, words: words })
            }
        }
    }

    return out
}

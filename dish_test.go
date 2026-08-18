package main
import("fmt"; "testing")

func TestWham(t *testing.T) {
    test(t, "tests/wham.dish")
}

func TestBoom(t *testing.T) {
    test(t, "tests/boom.dish")
}

func TestWhiz(t *testing.T) {
    test(t, "tests/whiz.dish")
}

func TestBang(t *testing.T) {
    test(t, "tests/bang.dish")
}

func TestWiki(t *testing.T) {
    test(t, "tests/wiki.dish")
}

func TestWaka(t *testing.T) {
    test(t, "tests/waka.dish")
}

func TestSplit(t *testing.T) {
    test(t, "tests/split.dish")
}

func TestSplat(t *testing.T) {
    test(t, "tests/splat.dish")
}

func TestDot(t *testing.T) {
    test(t, "tests/dot.dish")
}

func TestDash(t *testing.T) {
    test(t, "tests/dash.dish")
}

func TestSnip(t *testing.T) {
    test(t, "tests/snip.dish")
}

func TestSnap(t *testing.T) {
    test(t, "tests/snap.dish")
}

func TestTwiddle(t *testing.T) {
    test(t, "tests/twiddle.dish")
}

func TestGrep(t *testing.T) {
    test(t, "tests/grep.dish")
}

func TestPow(t *testing.T) {
    test(t, "tests/pow.dish")
}

func TestSharp(t *testing.T) {
    test(t, "tests/sharp.dish")
}

func TestCross(t *testing.T) {
    test(t, "tests/cross.dish")
}

func TestBars(t *testing.T) {
    test(t, "tests/bars.dish")
}

func TestThump(t *testing.T) {
    test(t, "tests/thump.dish")
}

func TestCriss(t *testing.T) {
    test(t, "tests/criss.dish")
}

func TestChirp(t *testing.T) {
    test(t, "tests/chirp.dish")
}

func TestBind(t *testing.T) {
    for _, c := range []struct{ in string; out string }{
        { in: "[0.1, 0.2]", out: "[0.1, 0.2]" },
        { in: "12345678901234567890", out: "12345678901234567890" },
        { in: "1.5e-3", out: "0.0015" },
        { in: "123 apples", out: "\"123 apples\"" },
        { in: "hello\nworld", out: "[\"hello\", \"world\"]" },
    } {
        if out := fmt.Sprintf("%v", bind([]byte(c.in))); out != c.out {
            t.Errorf("bind(%q) is %v; expected %v", c.in, out, c.out)
        }
    }
}

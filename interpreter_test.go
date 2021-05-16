package main
import "testing"

// Member ops
func TestMember(t *testing.T) {
    test(t, "ops/member.dish")
}

// Unary ops
func TestNot(t *testing.T) {
    test(t, "ops/not.dish")
}

func TestInverse(t *testing.T) {
    test(t, "ops/invert.dish")
}

// Binary ops
func TestTimes(t *testing.T) {
    test(t, "ops/multiply.dish")
}

func TestPlus(t *testing.T) {
    test(t, "ops/add.dish")
}

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

func TestProduct(t *testing.T) {
    test(t, "ops/product.dish")
}

func TestSum(t *testing.T) {
    test(t, "ops/sum.dish")
}

func TestNegate(t *testing.T) {
    test(t, "ops/negate.dish")
}

func TestStringify(t *testing.T) {
    test(t, "ops/str.dish")
}

func TestMin(t *testing.T) {
    test(t, "ops/min.dish")
}

func TestMax(t *testing.T) {
    test(t, "ops/max.dish")
}

func TestUnique(t *testing.T) {
    test(t, "ops/unique.dish")
}

func TestLength(t *testing.T) {
    test(t, "ops/length.dish")
}

func TestIncrement(t *testing.T) {
    test(t, "ops/increment.dish")
}

func TestDecrement(t *testing.T) {
    test(t, "ops/decrement.dish")
}

// Binary ops
func TestTimes(t *testing.T) {
    test(t, "ops/multiply.dish")
}

func TestPlus(t *testing.T) {
    test(t, "ops/add.dish")
}

func TestSubtract(t *testing.T) {
    test(t, "ops/subtract.dish")
}

func TestJoin(t *testing.T) {
    test(t, "ops/join.dish")
}

func TestBelow(t *testing.T) {
    test(t, "ops/below.dish")
}

func TestUnder(t *testing.T) {
    test(t, "ops/under.dish")
}

func TestAbove(t *testing.T) {
    test(t, "ops/above.dish")
}

func TestOver(t *testing.T) {
    test(t, "ops/over.dish")
}

func TestTo(t *testing.T) {
    test(t, "ops/to.dish")
}

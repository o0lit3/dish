package main
import "testing"

// Member ops
func TestMember(t *testing.T) {
    test(t, "tests/member.dish")
}

// Unary ops
func TestNot(t *testing.T) {
    test(t, "tests/not.dish")
}

func TestInverse(t *testing.T) {
    test(t, "tests/invert.dish")
}

func TestProduct(t *testing.T) {
    test(t, "tests/product.dish")
}

func TestKeys(t *testing.T) {
    test(t, "tests/keys.dish")
}

func TestItemize(t *testing.T) {
    test(t, "tests/itemize.dish")
}

func TestSum(t *testing.T) {
    test(t, "tests/sum.dish")
}

func TestNegate(t *testing.T) {
    test(t, "tests/negate.dish")
}

func TestPop(t *testing.T) {
    test(t, "tests/pop.dish")
}

func TestShift(t *testing.T) {
    test(t, "tests/shift.dish")
}

func TestStringify(t *testing.T) {
    test(t, "tests/stringify.dish")
}

func TestMin(t *testing.T) {
    test(t, "tests/min.dish")
}

func TestAverage(t *testing.T) {
    test(t, "tests/average.dish")
}

func TestMax(t *testing.T) {
    test(t, "tests/max.dish")
}

func TestUnique(t *testing.T) {
    test(t, "tests/uniq.dish")
}

func TestLength(t *testing.T) {
    test(t, "tests/length.dish")
}

func TestIncrement(t *testing.T) {
    test(t, "tests/increment.dish")
}

func TestDecrement(t *testing.T) {
    test(t, "tests/decrement.dish")
}

// Binary ops
func TestSwitch(t *testing.T) {
    test(t, "tests/switch.dish")
}

func TestFind(t *testing.T) {
    test(t, "tests/find.dish")
}

func TestPower(t *testing.T) {
    test(t, "tests/power.dish")
}

func TestMultiply(t *testing.T) {
    test(t, "tests/multiply.dish")
}

func TestDivide(t *testing.T) {
    test(t, "tests/divide.dish")
}

func TestMod(t *testing.T) {
    test(t, "tests/mod.dish")
}

func TestPlus(t *testing.T) {
    test(t, "tests/add.dish")
}

func TestSubtract(t *testing.T) {
    test(t, "tests/subtract.dish")
}

func TestJoin(t *testing.T) {
    test(t, "tests/join.dish")
}

func TestBase(t *testing.T) {
    test(t, "tests/base.dish")
}

func TestConvert(t *testing.T) {
    test(t, "tests/convert.dish")
}

func TestBelow(t *testing.T) {
    test(t, "tests/below.dish")
}

func TestUnder(t *testing.T) {
    test(t, "tests/under.dish")
}

func TestAbove(t *testing.T) {
    test(t, "tests/above.dish")
}

func TestOver(t *testing.T) {
    test(t, "tests/over.dish")
}

func TestEquals(t *testing.T) {
    test(t, "tests/equals.dish")
}

func TestAnd(t *testing.T) {
    test(t, "tests/and.dish")
}

func TestXor(t *testing.T) {
    test(t, "tests/xor.dish")
}

func TestOr(t *testing.T) {
    test(t, "tests/or.dish")
}

func TestRange(t *testing.T) {
    test(t, "tests/range.dish")
}

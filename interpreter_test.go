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

func TestHashify(t *testing.T) {
    test(t, "tests/hashify.dish")
}

func TestKeys(t *testing.T) {
    test(t, "tests/keys.dish")
}

func TestSort(t *testing.T) {
    test(t, "tests/sort.dish")
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

func TestCompact(t *testing.T) {
    test(t, "tests/compact.dish")
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

func TestRedo(t *testing.T) {
    test(t, "tests/redo.dish")
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

func TestRemainder(t *testing.T) {
    test(t, "tests/remainder.dish")
}

func TestAdd(t *testing.T) {
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

func TestPush(t *testing.T) {
    test(t, "tests/push.dish")
}

func TestUnshift(t *testing.T) {
    test(t, "tests/unshift.dish")
}

func TestIntersect(t *testing.T) {
    test(t, "tests/intersect.dish")
}

func TestExclude(t *testing.T) {
    test(t, "tests/exclude.dish")
}

func TestUnion(t *testing.T) {
    test(t, "tests/union.dish")
}

func TestLt(t *testing.T) {
    test(t, "tests/lt.dish")
}

func TestLte(t *testing.T) {
    test(t, "tests/lte.dish")
}

func TestGt(t *testing.T) {
    test(t, "tests/gt.dish")
}

func TestOver(t *testing.T) {
    test(t, "tests/gte.dish")
}

func TestCompare(t *testing.T) {
    test(t, "tests/compare.dish")
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

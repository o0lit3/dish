# Operators
Symbolic operators in **dish** (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods. The tables in this document outline which shorthand, symbolic operands correspond to which method names based on the data types of each operand.

**dish** is a (mostly) strongly-typed language. However, all data types in **dish** have inherit truthiness in Boolean contexts where the default values for each data type (`null`, Boolean `false`, Numeric `0`, String `""`, Array `()`, and Hash `{}`) evaluate to `false`, and all other values evaluate to `true`.

In addition, the Booleans `true` and `false` are treated as `1` and `0` in Numeric context, and `null` is treated as `0` in Numeric context. In the documentation below, the word "Number" can be replaced by a Boolean or Null coerced to a Number in this fashion. Any other type pairings that do not exist explicitly in this documentation will result in a runtime error (e.g. comparing a string to a number, `'3' > 2`).

With the exception of special [assignment operators](#assignment-operators) and [comparison operators](#comparison-operators) that use the `=` symbol as its second character, all symbolic operators in **dish** consist of either a single symbolic character or a doubled, symbolic character; `|`, `||`, and `|=` are a valid symbolic operators, `?:`, `=~`, and `>>=` are not.

## Logical Operators
Logical operators in **dish**, like many languages, short circuit the right-hand operand if the left-hand operand does not suffice to determine the value of the expression. As such, the right-hand operands of logical operators are implicit Blocks that are only evaluated if needed (or in the case of the `switch` and `swap` operators, an implicit Array of Blocks for both operands). Left-hand operands for logical operators are treated as implicit Booleans to determine truthiness, but are evaluated as the object itself when returned.

Logical operators and methods in **dish** serve as control structures for the language. **Nota bene**: Be careful with logical chains such as `(condition).then(0).else(1)`--falsy return values like 0 force the `else` block into execution regardless of the condition's truthiness. Therefore, to avoid any traps of falsy values, the ternary `(condition).switch[0, 1]` is preferred.

Similarly, to create traditional if/elseif/else logic in **dish**, use the n-ary `[cond1, cond2].switch[0, 1, 2]`. The conditions Array is executed in sequence until it finds a truthy value (with all subsequent conditions being short-circuited), and only the implicit Logic block from the right-hand Array corresponding to the same truthy index will be executed and returned. When no truthy condition is found, if the right-hand logic Array has more elements than the left-hand conditions Array, then the last index of the logic Array is executed and returned; otherwise `null` is returned.

| Operator | Operands              | Method Name         | Example                     | Result                |
| -------- | --------------------- | ------------------- | --------------------------- | --------------------- |
| `&&`     | Any `&&` Any          | `a.and(b)`          | `3 && 2`                    | `2`                   |
|          |                       |                     | `3 && 0`                    | `0`                   |
|          |                       |                     | `[] && 2`                   | `[]`                  |
|          |                       |                     |                             |                       |
| `\|\|`   | Any `\|\|` Any        | `a.or(b)`           | `3 \|\| 2`                  | `3`                   |
|          |                       |                     | `[] \|\| 0`                 | `0`                   |
|          |                       |                     | `0 \|\| 2`                  | `2`                   |
|          |                       |                     |                             |                       |
| `~~`     | Any `~~` Any          | `a.xor(b)`          | `3 ~~ 0`                    | `3`                   |
|          |                       |                     | `[] ~~ 2`                   | `2`                   |
|          |                       |                     | `3 ~~ 2`                    | `null`                |
|          |                       |                     |                             |                       |
| `??`     | Any `??` Any          | `a.coalesce(b)`     | `a ?? 2`                    | `2`                   |
|          |                       |                     | `a = 1; a ?? 2`             | `1`                   | 
|          |                       |                     |                             |                       |
| `?`      | Any `?` Number        | `a.then(b)`         | `(1 > 0) ? 2`               | `2`                   |
|          | Any `?` String        | `a.then(b)`         | `(0 > 1) ? 'binary'`        | `null`                |
|          | Any `?` Array         | `a.switch[b]`       | `[0, 1] ? [a += 1, a += 2]` | `2`                   |
|          | Any `?` :(Block)      | `a.while:(...)`     | `(a < 9) ? :(++a)`          | `9`                   |
|          |                       |                     |                             |                       |
| `!`      | Any `!` Number        | `a.else(b)`         | `(1 > 0) ! 2`               | `null`                |
|          | Any `!` String        | `a.else(b)`         | `(0 > 1) ! 'binary'`        | `"binary"`            |
|          | Any `!` Array         | `a.swap[b]`         | `[0, 1] ! [a += 1, a += 2]` | `1`                   |
|          | Any `!` :(Block)      | `a.until:(...)`     | `(a > 9) ! :(++a)`          | `10`                  |

## Comparison Operators
The `==` and `!=` equality operators in **dish** evaluate to `false` and `true` respectively when there is a type mismatch (including Booleans and Nulls that are otherwise coerced to `1` and `0` in Numeric contexts); all other comparison operators throw an error when there is a type mismatch. When operand types match, Numbers are compared numerically, Strings lexically, and Arrays and Hashes based on their string-evaluated value for equality (where hash keys are sorted alphabetically when string-evaluated), and by their length for inequalities.

| Operator | Operands                    | Method Name       | Example                        | Result            |
| -------- | --------------------------- | ----------------- | ------------------------------ | ----------------- |
| `==`     | Boolean `==` Boolean        | `a.eq(b)`         | `true == false`                | `false`           |
|          | Number `==` Number          | `a.eq(b)`         | `2 == 2.0`                     | `true`            |
|          | String `==` String          | `a.eq(b)`         | `'bin' == 'ary'`               | `false`           |
|          | Array `==` Array            | `a.is(b)`         | `[1, 2, 3] == [3, 2, 1]`       | `false`           |
|          | Hash `==` Hash              | `a.is(b)`         | `{a: 1, b: 2} == {b: 2, a: 1}` | `true`            |
|          |                             |                   |                                |                   |
| `!=`     | Boolean `!=` Boolean        | `a.ne(b)`         | `true != false`                | `true`            |
|          | Number `!=` Number          | `a.ne(b)`         | `2 != 2.0`                     | `false`           |
|          | String `!=` String          | `a.ne(b)`         | `'bin' != 'ary'`               | `true`            |
|          | Array `!=` Array            | `a.isnt(b)`       | `[1, 2, 3] != [3, 2, 1]`       | `true`            |
|          | Hash `!=` Hash              | `a.isnt(b)`       | `{a: 1, b: 2} != {b: 2, a: 1}` | `false`           |
|          |                             |                   |                                |                   |
| `>`      | Number `>` Number           | `a.gt(b)`         | `10 > 3.0`                     | `true`            |
|          | String `>` String           | `a.gt(b)`         | `'bin' > 'ary'`                | `true`            |
|          | Array `>` Array             | `a.gt(b)`         | `[4, 5] > [1, 2, 3]`           | `false`           |
|          | Hash `>` Hash               | `a.gt(b)`         | `{a: 1} > {b: 2, a: 1}`        | `false`           |
|          |                             |                   |                                |                   |
| `<`      | Number `<` Number           | `a.lt(b)`         | `10 < 3.0`                     | `false`           |
|          | String `<` String           | `a.lt(b)`         | `'bin' < 'ary'`                | `false`           |
|          | Array `<` Array             | `a.lt(b)`         | `[4, 5] < [1, 2, 3]`           | `true`            |
|          | Hash `<` Hash               | `a.lt(b)`         | `{a: 1} < {b: 2, a: 1}`        | `true`            |
|          |                             |                   |                                |                   |
| `>=`     | Number `>=` Number          | `a.ge(b)`         | `10 >= 3.0`                    | `true`            |
|          | String `>=` String          | `a.ge(b)`         | `'bin' >= 'ary'`               | `true`            |
|          | Array `>=` Array            | `a.ge(b)`         | `[4, 5] >= [1, 2]`             | `true`            |
|          | Hash `>=` Hash              | `a.ge(b)`         | `{a: 1} >= {b: 2, a: 1}`       | `false`           |
|          |                             |                   |                                |                   |
| `<=`     | Number `<=` Number          | `a.le(b)`         | `10 <= 3.0`                    | `false`           |
|          | String `<=` String          | `a.le(b)`         | `'bin' <= 'ary'`               | `false`           |
|          | Array `<=` Array            | `a.le(b)`         | `[4, 5] <= [1, 2]`             | `true`            |
|          | Hash `<=` Hash              | `a.le(b)`         | `{a: 1} <= {b: 2, a: 1}`       | `true`            |

## Assignment Operators
When initializing a variable for local scope, use `:`; when initializing a variable for global scope, use `=`. After initialization, `=` (and all variants) should be used to re-assign existing global _or_ local variables.

The member assignment operator `@=` is a special assignment operator that allows value assignment to a variable member (a Numeric bit, a String character, an Array index, or a Hash property), while returning the parent object.

| Operator | Operands                      | Method Name      | Example                   | Result                |
| -------- | ----------------------------- | ---------------- | ------------------------- | --------------------- |
| `:`      | Variable `:` Any              | `a.define(b)`    | `a: 5`                    | `5`                   |
|          |                               |                  |                           |                       |
| `:`      | `:`Variable                   |                  | `:a(...)`                 | `// Block argument`   |
| (Unary)  |                               |                  |                           |                       |
|          |                               |                  |                           |                       |
| `=`      | Variable `=` Any              | `a.assign(b)`    | `a = [1, 2, 3]`           | `[1, 2, 3]`           |
| `+=`     | [See Binary Ops](#binary-operators)              | | `a = 1, a += 1`         | `2`                   |
| `-=`     | [See Binary Ops](#binary-operators)              | | `a = 1, a -= 1`         | `0`                   |
| `*=`     | [See Binary Ops](#binary-operators)              | | `a = 2, a *= 2`         | `4`                   |
| `/=`     | [See Binary Ops](#binary-operators)              | | `a = 6, a /= 2`         | `3`                   |
| `%=`     | [See Binary Ops](#binary-operators)              | | `a = 5, a %= 4`         | `1`                   |
| `^=`     | [See Binary Ops](#binary-operators)              | | `a = 5, a ^= 3`         | `125`                 |
| `&=`     | [See Binary Ops](#binary-operators)              | | `a = 5, a &= 3`         | `1`                   |
| `~=`     | [See Binary Ops](#binary-operators)              | | `a = 5, a ~= 3`         | `6`                   |
| `\|=`    | [See Binary Ops](#binary-operators)              | | `a = 5, a \|= 3`        | `7`                   |
| `?=`     | [See Logical ?? Op](#logical-operators)          | | `a ?= 3`                | `3`                   |
|          |                               |                  |                           |                       |
| `@=`     | Number.x `@=` Implied Boolean | `a.x.replace(b)` | `a = 12, a.2 @= 0`        | `8`                   |
|          | String.x `@=` String          | `a.x.replace(b)` | `a = 'bing', a.1 @= 'a'`  | `"bang"`              |
|          | Array.x `@=` Any              | `a.x.replace(b)` | `a = [1, 2, 3], a.0 @= 4` | `[4, 2, 3]`           |
|          | Hash.x `@=` Any               | `a.x.replace(b)` | `a = {x: 1}, a.'x' @= 4`  | `{"x": 4}`            |
|          |                               |                  |                           |                       |
| `++`     | `++`Number                    | `a.increment`    | `a = 8, ++a`              | `9`                   |
| (Unary)  | `++`String                    | `a.increment`    | `a = 'binary', ++a`       | `"binarz"`            |
|          | `++`Array                     | `a.increment`    | `a = [1, 2, 3], ++a`      | `[1, 2, 3, null]`     |
|          |                               |                  |                           |                       |
| `--`     | `--`Number                    | `a.decrement`    | `a = 9, --a`              | `8`                   |
| (Unary)  | `--`String                    | `a.decrement`    | `a = 'binary', --a`       | `"binarx"`            |
|          | `--`Array                     | `a.decrement`    | `a = [1, 2, 3], --a`      | `[1, 2]`              |
|          |                               |                  |                           |                       |
| `<<`     | Null `<<` Any                 | `a.push(b)`      | `a << 'dish'`             | `["dish"]`            |
|          | String `<<` Any               | `a.append(b)`    | `'dish.' << 1.5`          | `"dish.1.5"`          |
|          | Array `<<` Any                | `a.push(b)`      | `[1, 2] << 3`             | `[1, 2, 3]`           |
|          | Hash `<<` Any                 | `a.extend(b)`    | `{x: 1} << {y: 2}`        | `{"x": 1, "y": 2}`    |
|          | [Also See Bitwise Ops](#bitwise-operators)       |                           |                       |
|          |                               |                  |                           |                       |
| `>>`     | Null `>>` Any                 | `a.unshift(b)`   | `a >> 'dish'`             | `["dish"]`            |
|          | String `>>` Any               | `a.prepend(b)`   | `'ary' >> 'bin'`          | `"binary"`            |
|          | Array `>>` Any                | `a.unshift(b)`   | `[1, 2] >> 3`             | `[3, 2, 1]`           |
|          | Hash `>>` Any                 | `a.extend(b)`    | `{x: 1} >> {y: 2}`        | `{"x": 1, "y": 2}`    |
|          | [Also See Bitwise Ops](#bitwise-operators)       |                           |                       |
|          |                               |                  |                           |                       |
| `<<`     | `<<`Number                    | `a.shift`        | `a = 12, [<<a, a]`        | `[1, 4]`              |
| (Unary)  | `<<`String                    | `a.shift`        | `a = 'binary', [<<a, a]`  | `["b", "inary"]`      |
|          | `<<`Array                     | `a.shift`        | `a = [1, 2, 3], [<<a, a]` | `[1, [2, 3]]`         |
|          | [Also See Unary Ops](#unary-operators)           |                           |                       |
|          |                               |                  |                           |                       |
| `>>`     | `>>`Number                    | `a.pop`          | `a = 12, [>>a, a]`        | `[0, 6]`              |
| (Unary)  | `>>`String                    | `a.pop`          | `a = 'binary', [>>a, a]`  | `["y", "binar"]`      |
|          | `>>`Array                     | `a.pop`          | `a = [1, 2, 3], [>>a, a]` | `[3, [1, 2]]`         |
|          | [Also See Unary Ops](#unary-operators)           |                           |                       |

**Dish** also supports parallel assignment when the left-hand operand is an Array of variables, as in `[a, b, c] = [1, 2, 3]`. Parallel assigment is useful for swapping the values held in two variables without the need of a temporary holding variable: `[a, b] = [b, a]` is equivalent to `t = a, a = b, b = t`.

This technique also works for member subsets, as in `a = [1, 2, 3]; a.[1, 2] = a.[2, 1]; a` or `a = [1, 2, 3]; a.[1, 2] @= a.[2, 1]`.

## Bitwise Operators
| Operator | Operands                 | Method Name         | Example                     | Result                 |
| -------- | ------------------------ | ------------------- | --------------------------- | ---------------------- |
| `&`      | Number `&` Number        | `a.band(b)`         | `5 & 3`                     | `1`                    |
|          | String `&` String        | `a.intersection(b)` | `'codfish' & 'airshed'`     | `"dish"`               |
|          | String `&` :(Block)      | `a.all:(...)`       | `'binary' & :c(c.vowel)`    | `false`                |
|          | Array `&` Array          | `a.intersection(b)` | `[1, 2, 3] & [4, 3, 2]`     | `[2, 3]`               |
|          | Array `&` :(Block)       | `a.all:(...)`       | `[1, 3, 5] & :x(x % 2)`     | `true`                 |
|          | Hash `&` Hash            | `a.intersection(b)` | `{x: 1} & {x: 2, y: 2}`     | `{"x": 2}`             |
|          | Hash `&` :(Block)        | `a.all:(...)`       | `{x: 1, y: 2} & :x(x % 2)`  | `false`                |
|          |                          |                     |                             |                        |
| `\|`     | Number `\|` Number       | `a.bor(b)`          | `5 \| 3`                    | `7`                    |
|          | String `\|` String       | `a.union(b)`        | `'didi' \| 'sihi'`          | `"dish"`               |
|          | String `\|` :(Block)     | `a.any:(...)`       | `'binary' \| :c(c.vowel)`   | `true`                 |
|          | Array `\|` Array         | `a.union(b)`        | `[1, 2, 3] \| [4, 2, 1]`    | `[1, 2, 3, 4]`         |
|          | Array `\|` :(Block)      | `a.any:(...)`       | `[1, 2, 3] \| :x(x % 2)`    | `true`                 |
|          | Hash `\|` Hash           | `a.union(b)`        | `{x: 1} \| {x: 2}`          | `{"x": 2}`             |
|          | Hash `\|` :(Block)       | `a.any:(...)`       | `{x: 2, y: 4} \| :x(x % 2)` | `false`                |
|          |                          |                     |                             |                        |
| `~`      | Number `~` Number        | `a.bxor(b)`         | `5 ~ 3`                     | `6`                    |
|          | String `~` String        | `a.exclusion(b)`    | `'dentist' ~ 'tenth'`       | `"dish"`               |
|          | String `~` :(Block)      | `a.none(...)`       | `'rsvp' ~ :c(c.vowel)`      | `true`                 |
|          | Array `~` Array          | `a.exclusion(b)`    | `[1, 2, 3] ~ [4, 3, 2]`     | `[1, 4]`               |
|          | Array `~` :(Block)       | `a.none(...)`       | `[2, 4, 6] ~ :x(x % 2)`     | `true`                 |
|          | Hash `~` Hash            | `a.exclusion(b)`    | `{x: 1} ~ {y: 2}`           | `{}`                   |
|          | Hash `~` :(Block)        | `a.none(...)`       | `{x: 1, y: 2} ~ :(x % 2)`   | `false`                |
|          |                          |                     |                             |                        |
| `~`      | `~`Number                | `a.bnot`            | `~5`                        | `-6`                   |
| (Unary)  | `~`String                | `a.caseflip`        | `~'bInArY'`                 | `"BiNaRy"`             |
|          | `~`Array                 | `a.transpose`       | `~[1, 2, 3]`                | `[[1], [2], [3]]`      |
|          | `~`Hash                  | `a.invert`          | `~{x: 1, y: 2}`             | `{"1": "x", "2": "y"}` |
|          |                          |                     |                             |                        |
| `<<`     | Number `<<` Number       | `a.lshift(b)`       | `5 << 3`                    | `40`                   |
|          | [Also See Assignment Ops](#assignment-operators)                         | | |                        |
|          |                          |                     |                             |                        |
| `>>`     | Number `>>` Number       | `a.rshift(b) `      | `40 >> 3`                   | `5`                    |
|          | [Also See Assignment Ops](#assignment-operators)                         | | |                        |

## Binary Operators
| Operator | Operands             | Method Name         | Example                        | Result                     |
| -------- | -------------------- | ------------------- | ------------------------------ | -------------------------- |
| `+`      | Number `+` Number    | `a.add(b)`          | `10 + 20`                      | `30`                       |
|          | Number `+` String    | `b.increase(a)`     | `2 + 'binary'`                 | `"binasa"`                 |
|          | Number `+` Array     | `b.lpad(a)`         | `2 + [1, 2, 3]`                | `[null, null, 1, 2, 3]`    |
|          | `-`Number `+` Array  | `b.ltrunc(a)`       | `-2 + [1, 2, 3]`               | `[3]`                      |
|          | String `+` Number    | `a.increase(b)`     | `'binary' + 2`                 | `"binasa"`                 |
|          | String `+` String    | `a.concat(b)`       | `'bin' + 'ary'`                | `"binary"`                 |
|          | String `+` :(Block)  | `a.accumulate:(..)` | `'bin' + :x:y(x + y.ord)`      | `[98, 203, 313]`           |
|          | Array `+` Number     | `a.rpad(b)`         | `[1, 2, 3] + 2`                | `[1, 2, 3, null, null]`    |
|          | Array `+` `-`Number  | `a.rtrunc(b)`       | `[1, 2, 3] + -2`               | `[1]`                      |
|          | Array `+` Array      | `a.concat(b)`       | `[1, 2] + [2]`                 | `[1, 2, 2]`                |
|          | Array `+` :(Block)   | `a.accumulate:(..)` | `[2, 3] + :x:y((x ?? 1) * y)`  | `[2, 6]`                   |
|          | Hash `+` Hash        | `a.concat(b)`       | `{x: 1} + {y: 2}`              | `{"x": 1, "y": 2}`         |
|          | Hash `+` :(Block)    | `a.accumulate:(..)` | `{x: 2, y: 3} + :x:y(x + y)`   | `[2, 5]`                   |
|          |                      |                     |                                |                            |
| `-`      | Number `-` Number    | `a.subtract(b)`     | `10 - 20`                      | `-10`                      |
|          | String `-` Number    | `a.decrease(b)`     | `'binary' - 2`                 | `"binarw"`                 |
|          | String `-` String    | `a.remove(b)`       | `'binary' - 'ary'`             | `"bin"`                    |
|          | String `-` :(Block)  | `a.reduce:(...)`    | `'bin' - :x:y(x + y.ord)`      | `313`                      |
|          | Array `-` Array      | `a.remove(b)`       | `[1, 2, 2, 3, 4] - [2, 3]`     | `[1, 2, 4]`                |
|          | Array `-` :(Block)   | `a.reduce:v(b)`     | `[2, 3] - :x:y((x ?? 1) * y)`  | `6`                        |
|          | Hash `-` Hash        | `a.remove(b)`       | `{x: 1, y: 2} - {x}`           | `{"y": 2}`                 |
|          | Hash `-` :(Block)    | `a.reduce:(b)`      | `{x: 2, y: 3} - :x:y(x + y)`   | `5`                        |
|          |                      |                     |                                |                            |
| `*`      | Number `*` Number    | `a.multiply(b)`     | `10 * 20`                      | `200`                      |
|          | Number `*` String    | `b.repeat(a)`       | `5 * 'a'`                      | `"aaaaa"`                  |
|          | Number `*` Array     | `b.repeat(a)`       | `2 * [1, 2, 3]`                | `[1, 2, 3, 1, 2, 3]`       |
|          | Number `*` :(Block)  | `a.times(b)`        | `4 * :(++a)`                   | `[1, 2, 3, 4]`             |
|          | String `*` Number    | `a.repeat(b)`       | `'a' * 5`                      | `"aaaaa"`                  |
|          | String `*` String    | `a.join(b)`         | `'bin' * 'ary'`                | `"baryiaryn"`              |
|          | String `*` Array     | `b.join(a)`         | `' ' * [1, 2, 3]`              | `"1 2 3"`                  |
|          | String `*` :(Block)  | `a.map:(...)`       | `'binary' * :c(++c)`           | `"cjobsz"`                 |
|          | Array `*` Number     | `a.repeat(b)`       | `[1, 2, 3] * 2`                | `[1, 2, 3, 1, 2, 3]`       |
|          | Array `*` String     | `a.join(b)`         | `[1, 2, 3] * ' '`              | `"1 2 3"`                  |
|          | Array `*` Array      | `a.dot(b)`          | `[1, 2, 3] * [1, 2, 3]`        | `[1, 4, 9]`                |
|          | Array `*` :(Block)   | `a.map:(...)`       | `[1, 2, 3] * :x(x * 2)`        | `[1, 4, 6]`                |
|          | Hash `*` Hash        | `a.dot(b)`          | `{x: 1, y: 2} * {y: 4}`        | `{"y": 8}`                 |
|          | Hash `*` :(Block)    | `a.map:(...)`       | `{x: 1, y: 2} * :x(x * 2)`     | `{"x": 2, "y": 4}`         |
|          |                      |                     |                                |                            |
| `**`     | Number `**` Number   | `a.choose(b)`       | `10 ** 2`                      | `45`                       |
|          | String `**` Number   | `a.comb(b)`         | `'bin' ** 2`                   | `["bi", "bn", "in"]`       |
|          | Array `**` Number    | `a.comb(b)`         | `[1, 2, 3] ** 2`               | `[[1, 2], [1, 3], [2, 3]]` |
|          | Hash `**` Number     | `a.comb(b)`         | `{x: 1, y: 2} ** 2`            | `[[1, 2]]`                 |
|          |                      |                     |                                |                            |
| `/`      | Number `/` Number    | `a.divide(b)`       | `10 / 20`                      | `0.5`                      |
|          | String `/` Number    | `a.split(b)`        | `'binary' / 2`                 | `["bin", "ary"]`           |
|          | String `/` String    | `a.split(b)`        | `'binary' / 'n'`               | `["bi", "ary"]`            |
|          | String `/` :(Block)  | `a.split:(...)`     | `'binary' / :c(c.vowel)`       | `["b", "n", "ry"]`         |
|          | Array `/` Number     | `a.split(b)`        | `[1, 2, 3, 4, 5] / 2`          | `[[1, 2, 3], [4, 5]]`      |
|          | Array `/` :(Block)   | `a.split:(...)`     | `[2, 2, 3, 4] / :x(x % 2)`     | `[[2, 2], [4]]`            |
|          |                      |                     |                                |                            |
| `//`     | Number `//` Number   | `a.idiv(b)`         | `250 // 24`                    | `10`                       |
|          | String `//` Number   | `a.partition(b)`    | `'binary' // 2`                | `["bi", "na", "ry"]`       |
|          | String `//` String   | `a.partition(b)`    | `'binary' // 'n'`              | `["bin", "ary"]`           |
|          | String `//` :(Block) | `a.group:(...)`     | `'Egg' // :c(c.ord)`           | `{"103": "gg", "69": "e"}` |
|          | Arary `//` Number    | `a.partition(b)`    | `[1, 2, 3, 4, 5] // 2`         | `[[1, 2], [3, 4], [5]]`    |
|          | Array `//` :(Block)  | `a.group:(...)`     | `[2, 3, 4] // :x(x % 2)`       | `{"0": [2, 4], "1": [3]}`  |
|          |                      |                     |                                |                            |
| `%`      | Number `%` Number    | `a.mod(b)`          | `3.5 % 2`                      | `1.5`                      |
|          | String `%` Number    | `a.every(b)`        | `'binary' % 2`                 | `"bnr"`                    |
|          | String `%` String    | `a.filter(b)`       | `'badabing' % 'b'`             | `"bb"`                     |
|          | String `%` :(Block)  | `a.filter:(...)`    | `'binary' % :c(c.vowel)`       | `"ia"`                     |
|          | Array `%` Number     | `a.every(b)`        | `[1, 2, 3, 4, 5, 6, 7] % 3`    | `[1, 4, 7]`                |
|          | Array `%` :(Block)   | `a.filter:(...)`    | `[1, 1, 2, 3] % :x(x % 2)`     | `[1, 1, 3]`                |
|          | Hash `%` :(Block)    | `a.filter:(...)`    | `{x: 1, y: 2} % :x(x % 2)`     | `{x: 1}`                   |
|          |                      |                     |                                |                            |
| `%%`     | Number `%%` Number   | `a.imod(b)`         | `3.5 %% 2`                     | `1`                        |
|          | String `%%` Number   | `a.xevery(b)`       | `'binary' %% 2`                | `"iay"`                    |
|          | String `%%` String   | `a.without(b)`      | `'badabing' %% 'b'`            | `"adaing"`                 |
|          | String `%%` :(Block) | `a.without:(...)`   | `'binary' %% :c(c.vowel)`      | `"bnry"`                   |
|          | Array `%%` Number    | `a.xevery(b)`       | `[1, 2, 3, 4, 5, 6, 7] %% 3`   | `[2, 3, 5, 6]`             |
|          | Array `%%` Array     | `a.without(b)`      | `[1, 2, 2, 3, 4] %% [2, 3]`    | `[1, 4]`                   |
|          | Array `%%` :(Block)  | `a.without:(...)`   | `[1, 1, 2, 3] %% :x(x % 2)`    | `[2]`                      |
|          | Hash `%%` :(Block)   | `a.without:(...)`   | `{x: 1, y: 2} %% :x(x % 2)`    | `{x: 2}`                   |
|          |                      |                     |                                |                            |
| `^`      | Number `^` Number    | `a.pow(b)`          | `2 ^ 3`                        | `8`                        |
|          | String `^` Number    | `a.rotate(b)`       | `'binary' ^ 2`                 | `"rybina"`                 |
|          | String `^` String    | `a.zip(b)`          | `'ds' ^ 'ih`                   | `"dish"`                   |
|          | String `^` :(Block)  | `a.sort:(...)`      | `'binary' ^ :x:y(y < x)`       | `"yrniba"`                 |
|          | Array `^` Number     | `a.rotate(b)`       | `[7, 9, 4] ^ 1`                | `[4, 7, 9]`                |
|          | Array `^` Array      | `a.zip(b)`          | `[1, 2, 3] ^ [4, 5]`           | `[[1, 4], [2, 5], [3]]`    |
|          | Array `^` :(Block)   | `a.sort:(...)`      | `[7, 9, 4] ^ :x:y(y < x)`      | `[9, 7, 4]`                |
|          | Hash `^` Hash        | `a.zip(b)`          | `{x: 1} ^ {x: 3}`              | `{"x": [1, 3]}`            |
|          |                      |                     |                                |                            |
| `.`      | Number`.`(Number)    | `a.at(b)`           | `12.(2)`                       | `1`                        |
|          | Number`.`Array       | `a.subset(b)`       | `12.[0, 1, 2]`                 | `4`                        |
|          | Number`.`:(Block)    | `a.call:(...)`      | `12.:n(+@n)`                   | `2`                        |
|          | String`.`Number      | `a.at(b)`           | `'binary'.1`                   | `"i"`                      |
|          | String`.`Array       | `a.subset(b)`       | `'binary'.[0, 3, 4]`           | `"bar"`                    |
|          | String`.`:(Block)    | `a.call:(...)`      | `'binary'.:s(s.uc)`            | `"BINARY`                  |
|          | Array`.`Number       | `a.at(b)`           | `[1, 2, 3].2`                  | `3`                        |
|          | Array`.`Array        | `a.items(b)`        | `[1, 2, 3].[0, 2]`             | `[1, 3]`                   |
|          | Array`.`:(Block)     | `a.call:(...)`      | `[1, 2, 3].:a(*a)`             | `6`                        |
|          | Hash`.`String        | `a.at(b)`           | `{x: 1, y: 2}.'y'`             | `2`                        |
|          | Hash`.`Array         | `a.items(b)`        | `{x: 1, y: 2}.['x', 'y']`      | `[1, 2]`                   |
|          | Hash`.`:(Block)      | `a.call:(...)`      | `{x: 1, y: 2}.:h(+h)`          | `3`                        |
|          |                      |                     |                                |                            |
| `@`      | Number `@` Number    | `a.round(b)`        | `3.1415 @ 3`                   | `3.142`                    |
|          | String `@` String    | `a.find(b)`         | `'binary' @ 'ary'`             | `[3]`                      |
|          | String `@` :(Block)  | `a.search:(...)`    | `'binary' @ :c(c.vowel)`       | `[1, 3]`                   |
|          | Array `@` Any        | `a.find(b)`         | `[1, 3, 2, 3] @ 3`             | `[1, 3]`                   |
|          | Array `@` :(Block)   | `a.search:(...)`    | `[1, 2, 3] @ :n(n % 2)`        | `[0, 2]`                   |
|          | Hash `@` Any         | `a.find(b)`         | `{a: 1, b: 2, c: 2} @ 2`       | `["b", "c"]`               |
|          | Hash `@` :(Block)    | `a.search:(...)`    | `{f: 1, b: 2} @ :n(n % 2)`     | `["f"]`                    |
|          |                      |                     |                                |                            |
| `#`      | Null `#` String      | `a.fmt(b)`          | `a # 'a: %s'`                  | `"a: null"`                |
|          | Boolean `#` Number   | `a.fmt(b)`          | `a = true, a # 'a: %s'`        | `"a: true"`                |
|          | Number `#` Number    | `a.base(b)`         | `9 # 2`                        | `"1001"`                   |
|          | Number `#` String    | `a.fmt(b)`          | `9 # '%.2f'`                   | `"9.00"`                   |
|          | String `#` Number    | `a.unbase(b)`       | `'1001' # 2`                   | `9`                        |
|          | String `#` String    | `a.fmt(b)`          | `'binary' # '%7s'`             | `" binary"`                |
|          | Array `#` String     | `a.fmt(b)`          | `[1, 2, 3] # '%s, %s, and %s'` | `"1, 2, and 3"`            |
|          | Hash `#` String      | `a.fmt(b)`          | `{x: 1, y: 2} # '%s and %s'`   | `"1 and 2"`                |
|          |                      |                     |                                |                            |
| `..`     | Number `..` Number   | `a.to(b)`           | `3 .. 1`                       | `[3, 2, 1]`                |
|          | String `..` String   | `a.to(b)`           | `'a' .. 'c'`                   | `["a", "b", "c"]`          |

**Nota bene**: the `..` range/to operator supports positve as well as negative ranges (as forced integers for numbers). In cases where a range is used to short-circuit a traditional for loop, you may need to first check that you haven't entered a negative range. Compare the following:

`dish -e 'a = 1; (0..3).map:i((1..i).map:j(a *= j)); a'` outputs `0`

`dish -e 'a = 1; (0..3).map:i((i > 1).then((1..i).map:j(a *= j))); a'` outputs `12`

The `..` range/to operator can also operate on strings, applying the string increment or decrement operation on the left-hand operand until the right-hand operand has been reached, or until the length of the right-hand operand has been exhausted (be careful with memory consumption when writing such operations). The String increment and decrement operators work on three mutally exclusive ranges ('0'..'9'`, `'A'..'Z'`, and `'a'..'z'`), using "carry over" addition/subtraction on the preceding character when the current character range has been eclipsed. For example:

`dish -e '["z" + 1, "z" - 1, "z" + 27]'` outputs `["aa", "y", "ba"]`

`dish -e '["9" + 1, "9" - 1, "9" + 11]'` outputs `["10", "8", "20"]`

`dish -e '["9z" + 1, "9z" - 1, "9z" + 11, "9z" + 27]'` outputs `["10a", "9y", "10k", "11a"]`

All characters that are outside the above mentioned character ranges are ignored on String increment and decrement operations. Note from the second example above that increasing or decreasing a numeric String has a similar effect of performing the same operation on an actual Number, but returns a String.

## Unary Operators
| Operator | Operands                | Method Name   | Example                  | Result                         |
| -------- | ----------------------- | ------------- | ------------------------ | ------------------------------ |
| `!`      | `!`Implied Boolean      | `a.not`       | `!5`                     | `false`                        |
|          |                         |               |                          |                                |
| `?`      | `?`Implied Boolean      | `a.bool`      | `?5`                     | `true`                         |
|          |                         |               |                          |                                |
| `+`      | `+`Number               | `a.num`       | `+5`                     | `5`                            |
|          | `+`String               | `a.num`       | `+"5"`                   | `5`                            |
|          | `+`Array                | `a.sum`       | `+[1, 2, 3]`             | `6`                            |
|          | `+`Hash                 | `a.sum`       | `+{x: 2, y: 3}`          | `5`                            |
|          |                         |               |                          |                                |
| `++`     | [See Assignment Ops](#assignment-operators)                    | | |                                |
|          |                         |               |                          |                                |
| `-`      | `-`Number               | `a.negate`    | `-5`                     | `-5`                           |
|          | `-`String               | `a.separate`  | `-"bin"`                 | `["b", "i", "n"]`              |
|          | `-`Array                | `a.negsum`    | `-[1, 2, 3]`             | `-6`                           |
|          | `-`Hash                 | `a.negsum`    | `-{x: 2, y: 3}`          | `-5`                           |
|          |                         |               |                          |                                |
| `--`     | [See Assignment Ops](#assignment-operators)                    | | |                                |
|          |                         |               |                          |                                |
| `*`      | `*`Null                 | `a.str`       | `*a`                     | `"null"`                       |
|          | `*`Boolean              | `a.str`       | `*true`                  | `"true"`                       |
|          | `*`Number               | `a.str`       | `*7`                     | `"7"`                          |
|          | `*`String               | `a.str`       | `*'binary'`              | `"binary"`                     |
|          | `*`Array                | `a.prod`      | `*[1, 2, 3]`             | `6`                            |
|          | `*`Hash                 | `a.prod`      | `*{x: 2, y: 3}`          | `6`                            |
|          |                         |               |                          |                                |
| `**`     | `**`Number              | `a.divisors`  | `**12`                   | `[1, 2, 3, 4, 6]`              |
|          | `**`String              | `a.perms`     | `**'ab'`                 | `["ab", "ba"]`                 |
|          | `**`Array               | `a.perms`     | `**[1, 2]`               | `[[1, 2], [2, 1]]`             |
|          | `**`Hash                | `a.perms`     | `**{x: 1, y: 2}`         | `[{x: 1, y: 2}, {x: 2, y: 1}]` |
|          |                         |               |                          |                                |
| `/`      | `/`Number               | `a.factors`   | `/12`                    | `[2, 2, 3]`                    |
|          | `/`String               | `a.split`     | `/'Hello World!'`        | `["Hello", "World!"]`          |
|          | `/`Array                | `a.array`     | `/[1, 2, 3]`             | `[1, 2, 3]`                    |
|          | `/`Hash                 | `a.vals`      | `/{x: 1, y: 2}`          | `[1, 2]`                       |
|          |                         |               |                          |                                |
| `//`     | `//`Number              | `a.bits`      | `//12`                   | `[0, 0, 1, 1]`                 |
|          | `//`String              | `a.chars`     | `//'bin'`                | `["b", "i", "n"]`              |
|          | `//`Array               | `a.flat`      | `//[1, [2, 3], 4]`       | `[1, 2, 3, 4]`                 |
|          | `//`Hash                | `a.flat`      | `//{x: 1, y: 2}`         | `["x", 1, "y", 2]`             |
|          |                         |               |                          |                                |
| `%`      | `%`Number               | `a.ratio`     | `%5`                     | `{"num": 157, "denom": 50}`    |
|          | `%`String               | `a.hash`      | `%'x'`                   | `{"x": true}`                  |
|          | `%`Array                | `a.hash`      | `%[2]`                   | `{"2": true}`                  |
|          | `%`Hash                 | `a.hash`      | `%{x: 1, y: 2}`          | `{"x": 1, "y": 2}`             |
|          |                         |               |                          |                                |
| `@`      | `@`Number               | `a.round`     | `@3.14`                  | `3`                            |
|          | `@`String               | `a.reverse`   | `@'binary'`              | `"yranib"`                     |
|          | `@`Array                | `a.reverse`   | `@[1, 2, 3]`             | `[3, 2, 1]`                    |
|          | `@`Hash                 | `a.keys`      | `@{x: 1, y: 2}`          | `["x", "y"]`                   |
|          |                         |               |                          |                                |
| `#`      | `#`Null                 | `a.len`       | `#null`                  | `0`                            |
|          | `#`Boolean              | `a.bitcount`  | `[#true, #false]`        | `[1, 0]`                       |
|          | `#`Number               | `a.bitcount`  | `[#12, #3, #2, #1, #0]`  | `[4, 2, 2, 1, 0]`              |
|          | `#`String               | `a.len`       | `#'binary'`              | `6`                            |
|          | `#`Array                | `a.len`       | `#[1, 2, 3]`             | `3`                            |
|          | `#`Hash                 | `a.len`       | `#{x: 1, y: 2}`          | `2`                            |
|          |                         |               |                          |                                |
| `<<`     | `<<`Number              | `a.msb`       | `<<12`                   | `1`                            |
|          | `<<`String              | `a.first`     | `<<'binary'`             | `"b"`                          |
|          | `<<`Array               | `a.first`     | `<<[1, 2, 3]`            | `1`                            |
|          | [Also See Assignment Ops](#assignment-operators)               | | |                                |
|          |                         |               |                          |                                |
| `>>`     | `>>`Number              | `a.lsb`       | `>>12`                   | `0`                            |
|          | `>>`String              | `a.last`      | `>>'binary'`             | `"y"`                          |
|          | `>>`Array               | `a.last`      | `>>[1, 2, 3]`            | `3`                            |
|          | [Also See Assignment Ops](#assignment-operators)               | | |                                |
|          |                         |               |                          |                                |
| `^`      | `^`Number               | `a.squared`   | `^9`                     | `81`                           |
|          | `^`String               | `a.sort`      | `^'binary'`              | `"abinry"`                     |
|          | `^`Array                | `a.sort`      | `^[10, 1, 4, 2]`         | `[1, 2, 4, 10]`                |
|          |                         |               |                          |                                |
| `<`      | `<`Number               | `a.floor`     | `<3.14`                  | `3`                            |
|          | `<`String               | `a.lc`        | `<'BINARY'`              | `"binary"`                     |
|          | `<`Array                | `a.min`       | `<[1, 2, 3]`             | `1`                            |
|          | `<`Hash                 | `a.min`       | `<{x: 1, y: 2}`          | `1      `                      |
|          |                         |               |                          |                                |
| `>`      | `>`Number               | `a.ceil`      | `>3.14`                  | `4`                            |
|          | `>`String               | `a.uc`        | `>'binary'`              | `"BINARY"`                     |
|          | `>`Array                | `a.max`       | `>[1, 2, 3]`             | `3`                            |
|          | `>`Hash                 | `a.max`       | `>{x: 1, y: 2}`          | `2`                            |
|          |                         |               |                          |                                |
| `=`      | `=`Number               | `a.num`       | `=3.14`                  | `3.14`                         |
|          | `=`String               | `a.eval`      | `='2 + 2'`               | `4`                            |
|          | `=`Array                | `a.avg`       | `=[1, 2, 3]`             | `2`                            |
|          | `=`Hash                 | `a.avg`       | `={x: 1, y: 2}`          | `1.5`                          |
|          |                         |               |                          |                                |
| `~`      | [See Bitwise Ops](#bitwise-operators) | |                          |                                |
|          |                         |               |                          |                                |
| `&`      | `&`Number               | `a.popcount   | `&12`                    | `2`                            |
|          | `&`String               | `a.compact`   | `&'Hello World'`         | `"HelloWorld`                  |
|          | `&`Array                | `a.compact`   | `&[null, 1, 2]`          | `[1, 2]`                       |
|          | `&`Hash                 | `a.compact`   | `&{a: 0, b: 1, c: null}` | `{"b": 1}`                     |
|          |                         |               |                          |                                |
| `\|`     | `\|`Number              | `a.abs`       | `\|-3`                   | `3`                            |
|          | `\|`String              | `a.uniq`      | `\|'bookkeeper'`         | `"bokepr"`                     |
|          | `\|`Array               | `a.uniq`      | `\|[1, 1, 2, 3, 3]`      | `[1, 2, 3]`                    |
|          | `\|`Hash                | `a.uniq`      | `\|{x: 1, y: 1, z: 2}`   | `{"x": 1, "z": 2}`             |

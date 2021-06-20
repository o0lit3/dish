# Operators
Symbolic operators in **dish** (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods. The tables below outline which shorthand, symbolic operands correspond to which method names. Some binary operators change behavior depending on the data types of their operands; in these cases, the left-hand operand forces the right-hand operand into an implicit context with the following conversion rules:

|             | Boolean               | Number      | String              | Array      | Hash      |
| ----------- | --------------------- | ----------- | ------------------- | ---------- | --------- |
| **Boolean** | `x`                   | `1`, `0`    | `"true"`, `"false"` | `[x]`      | `{~x: x}` |
| **Number**  | `x != 0`              | `x`         | `~x`                | `[x]`      | `{~x: x}` |
| **String**  | `x != "0" && x != ""` | `+x`        | `x`                 | `[x]`      | `{~x: x}` |
| **Array**   | `#x != 0`             | `#x`        | `x.join('')`        | `x`        | `x.hash`  |
| **Hash**    | `#x != 0`             | `#x`        | `x.values.join('')` | `x.values` | `x`       |

With the exception of special [assignment operators](#assignment-operators) and [comparison operators](#comparison-operators), all symbolic operators in **dish** consist of either a single symbolic character or a doubled, symbolic character; `|` and `||` are a valid symbolic operators, `?:` is not.

For the purposes of the tables below, "Scalar" refers to a Boolean or Number in numeric context; "List" refers to an Array or Hash in array context.

## Logical Operators
Logical operators in **dish**, like many languages, short circuit the right-hand argument if the left-hand argument does not suffice to determine the value of the expression. As such, the right-hand argument of logical operators are implicit Blocks that are only evaluated if needed (or in the case of the `switch` operator, an array of implicit Blocks for both operands). These operators return either `false` or a truthy value of any data type.

Logical operators and methods in **dish** serve as the control structures for the language.

| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `&&`     | Implicit Boolean && Any     | `a.then(b)`       | `3 && 2`                   | `2`                   |
|          |                             |                   | `true && 0`                | `false`               |
|          |                             |                   |                            |                       |
| `\|\|`   | Implicit Boolean || Any     | `a.else(b)`       | `3 || 2`                   | `3`                   |
|          |                             |                   | `false || 0`               | `false`               |
|          |                             |                   |                            |                       |
| `^^`     | Implicit Boolean ^^ Any     | `a.xor(b)`        | `3 ^^ 2`                   | `false`               |
|          |                             |                   | `false ^^ 2`               | `2`                   |
|          |                             |                   |                            |                       |
| `?`      | Boolean Array ? Any Array   | `a.switch[b]`     | `[1 > 0] ? [3, 4]`         | `3`                   |
|          |                             |                   | `[1 < 0] ? [3, 4]`         | `4`                   |
|          |                             |                   |                            |                       |
| `??`     | Implicit Boolean ?? Any     | `a.redo(b)`       | `(a < 9) ?? (++a)`         | `9`                   |
|          |                             |                   | `(1 < 0) ?? (++a)`         | `false`               |

## Comparison Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `==`     | Scalar == Implicit Number   | `a.equals(b)`     | `2 == 2`                   | `true`                |
|          | String == Implicit String   | `a.equals(b)`     | `"2" == "2.0"`             | `false`               |
|          | List == Implicit List       | `a.equals(b)`     | `[1, 2, 3] == [1, 2, 3]`   | `true`                |
|          |                             |                   |                            |                       |
| `!=`     | Scalar != Implicit Number   | `a.isnt(b)`       | `2 != 2`                   | `false`               |
|          | String != Implicit String   | `a.isnt(b)`       | `"2" != "2.0"`             | `true`                |
|          | List != Implicit List       | `a.isnt(b)`       | `[1, 2, 3] != [1, 2, 3]    | `false`               |
|          |                             |                   |                            |                       |
| `>`      | Scalar > Implicit Number    | `a.above(b)`      | `3 > 10`                   | `false`               |
|          | String > String             | `a.above(b)`      | `"3" > "10"`               | `true`                |
|          | String > Implicit Number    | `a.above(b)`      | `"3" > 10`                 | `false`               |
|          | List > Implicit Number      | `a.above(b)`      | `[1, 1, 1] > [2, 2]`       | `true`                |
|          |                             |                   |                            |                       |
| `<`      | Scalar < Implicit Number    | `a.below(b)`      | `3 < 10`                   | `true`                |
|          | String < String             | `a.below(b)`      | `"3" < "10"`               | `false`               |
|          | String < Implicit Number    | `a.below(b)`      | `"3" < 10`                 | `true`                |
|          | List < Implicit Number      | `a.below(b)`      | `[1, 1, 1] < [2, 2]`       | `false`               |
|          |                             |                   |                            |                       |
| `>=`     | Scalar >= Implicit Number   | `a.over(b)`       | `3 >= 10`                  | `false`               |
|          | String >= String            | `a.over(b)`       | `"3" >= "10"`              | `true`                |
|          | String >= Implicit Number   | `a.over(b)`       | `"3" >= 10`                | `false`               |
|          | List >= Implicit Number     | `a.over(b)`       | `[1, 1] >= [2, 2]`         | `true`                |
|          |                             |                   |                            |                       |
| `<=`     | Scalar <= Implicit Number   | `a.under(b)`      | `3 <= 10`                  | `true`                |
|          | String <= String            | `a.under(b)`      | `"3" <= "10"`              | `false`               |
|          | String <= Implicit Number   | `a.under(b)`      | `"3" <= 10`                | `true`                |
|          | List <= Implicit Number     | `a.under(b)`      | `[1, 1] <= [2, 2]`         | `true`                |

## Assignment Operators
When initializing a variable in global scope, use `=`; when initializing a variable in local scope, use `:`. After initialization, `=` (and all variants) should be used to re-assign existing global _or_ local variable.

| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `=`      | Variable = Any              | `a.assign(b)`     | `a = [1, 2, 3]`            | `[1, 2, 3]            |
| `~=`     |                             |                   | `a = 'bin', a ~= 'ary'     | `"binary"             |
| `+=`     |                             |                   | `a = 1, a += 1`            | `2`                   |
| `-=`     |                             |                   | `a = 1, a -= 1`            | `0`                   |
| `*=`     |                             |                   | `a = 2, a *= 2`            | `4`                   |
| `/=`     |                             |                   | `a = 6, a /= 2`            | `3`                   |
| `%=`     |                             |                   | `a = 5, a %= 4`            | `1`                   |
| `&=`     |                             |                   | `a = 5, a &= 3`            | `1`                   |
| `^=`     |                             |                   | `a = 5, a ^= 3`            | `6`                   |
| `\|=`    |                             |                   | `a = 5, a |= 3`            | `7`                   |

**Dish** also supports parallel assignment when the left-hand operand is an Array of variables, as in `[a, b, c] = [1, 2, 3]`. Parallel assigment is useful for swapping the values held in two variables without the need of a temporary holding variable: `[a, b] = [b, a]` is equivalent to `t = a, a = b, b = t`.

## Binary Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `~`      | Scalar + Implicit String    | `a.concat(b)`     | `10 ~ 20`                  | `"1020"`              |
|          | String + Implicit String    | `a.concat(b)`     | `"10" ~ "20"`              | `"1020"`              |
|          | List + Implicit String      | `a.join(b)`       | `[1, 2, 3] ~ ' '`          | `"1 2 3"`             |
|          |                             |                   |                            |                       |
| `+`      | Scalar + Implicit Number    | `a.add(b)`        | `10 + 20`                  | `30`                  |
|          | String + Implicit Number    | `a.add(b)`        | `'10' + '20'`              | `30`                  |
|          | Array + Implicit Array      | `a.add(b)`        | `[1, 2] + [2, 3]`          | `[1, 2, 2, 3]`        |
|          | Hash + Implicit Hash        | `a.add(b)`        | `{x: 1} + {y: 2}`          | `{"x": 1, "y": 2}`    |
|          |                             |                   |                            |                       |
| `-`      | Scalar - Implicit Number    | `a.subtract(b)`   | `10 - 20`                  | `-10`                 |
|          | String - Implicit Number    | `a.subtract(b)`   | `'10' - '20'`              | `-10`                 |
|          | Array - Implicit Array      | `a.remove(b)`     | `[1, 2] - [2, 3]`          | `[1]`                 |
|          | Hash - Implicit Hash        | `a.remove(b)`     | `{x: 1, y: 2} - {"x"}`     | `{"y": 2}`            |
|          |                             |                   |                            |                       |
| `*`      | Scalar * Implicit Number    | `a.multiply(b)`   | `10 * 20`                  | `200`                 |
|          | String * Implicit Number    | `a.repeat(b)`     | `"a" * 5`                  | `"aaaaa"`             |
|          | List * :(Block)             | `a.map:i(...)`    | `[1, 2, 3] * :i(i + 1)`    | `[2, 3, 4]`           |
|          | List * List                 | `a.dot(b)`        | `[1, 2] * [3, 4]`          | `[3, 8]`              |
|          | List * Implicit Number      | `a.repeat(b)`     | `[2] * 5`                  | `[2, 2, 2, 2, 2]`     |
|          |                             |                   |                            |                       |
| `/`      | Scalar / Implicit Number    | `a.divide(b)`     | `10 / 20`                  | `0.5`                 |
|          | String / :(Block)           | `a.split:i(...)`  | `'binary' / :i(i.vowel)`   | `["b", "n", "ry"]`    |
|          | String / String             | `a.split(b)`      | `'binary' / 'in'`          | `["b", "ary"]`        |
|          | String / Implicit Number    | `a.split(b)`      | `'binary' / 2`             | ` ["bi", "na", "ry"]` |
|          | List / :(Block)             | `a.split:i(...)`  | `[2, 2, 3, 4] / :i(i % 2)` | `[[2, 2], [4]]`       |
|          | List / Implicit Number      | `a.split(b)`      | `[2, 2, 3, 4] / 2`         | `[[2, 2], [3, 4]]`    |
|          |                             |                   |                            |                       |
| `%`      | Scalar % Implicit Number    | `a.remainder(b)`  | `20 % 7`                   | `6`                   |
|          | String % :(Block)           | `a.filter:i(...)` | `'binary' % :i(i.vowel)`   | `"ia"`                |
|          | String % Implicit Number    | `a.every(b)`      | `'binary' % 2`             | `"bnr"`               |
|          | List % :(Block)             | `a.filter:i(...)` | `[1, 1, 2, 3] % :i(i % 2)` | `[1, 1, 3]`           |
|          | List % Implicit Number      | `a.every(b)`      | `[1, 1, 2, 3] % 2`         | `[1, 2]`              |
|          |                             |                   |                            |                       |
| `**`     | Scalar \*\* Implicit Number | `a.pow(b)`        | `2 ** 3`                   | `8`                   |
|          | String \*\* :(Block)        | `a.sort:x:y(...)` | `'binary' ** :x:y(y < x)`  | `"yrniba"`            | 
|          | String \*\* Implicit Number | `a.pow(b)`        | `"2" ** "3"`               | `8`                   |
|          | List \*\* :(Block)          | `a.sort:x:y(...)` | `[7, 9, 4] ** :x:y(y < x)` | `[9, 7, 4]`           |
|          | List \*\* Implicit Number   | `a.rotate(b)`     | `[7, 9, 4] ** 1`           | `[4, 7, 9]`           |
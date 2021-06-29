# Operators
Symbolic operators in **dish** (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods. The tables below outline which shorthand, symbolic operands correspond to which method names. Some binary operators change behavior depending on the data types of their operands; in these cases, the left-hand operand forces the right-hand operand into an implicit context with the following conversion rules:

|             | Boolean               | Number      | String              | Array      | Hash      |
| ----------- | --------------------- | ----------- | ------------------- | ---------- | --------- |
| **Boolean** | `x`                   | `1`, `0`    | `"true"`, `"false"` | `[x]`      | `{~x: x}` |
| **Number**  | `x != 0`              | `x`         | `~x`                | `[x]`      | `{~x: x}` |
| **String**  | `x != "0" && x != ""` | `+x`        | `x`                 | `[x]`      | `{~x: x}` |
| **Array**   | `#x != 0`             | `#x`        | `x.join('')`        | `x`        | `x.hash`  |
| **Hash**    | `#x != 0`             | `#x`        | `x.values.join('')` | `x.values` | `x`       |

With the exception of special [assignment operators](#assignment-operators) and [comparison operators](#comparison-operators) that use the `=` symbol as its second character, all symbolic operators in **dish** consist of either a single symbolic character or a doubled, symbolic character; `|`, `||`, and `|=` are a valid symbolic operators, `?:`, `=~`, and `>>=` are not.

For the purposes of the tables below, "Scalar" refers to a Boolean or Number in numeric context; "List" refers to an Array or Hash in array context.

## Logical Operators
Logical operators in **dish**, like many languages, short circuit the right-hand argument if the left-hand argument does not suffice to determine the value of the expression. As such, the right-hand argument of logical operators are implicit Blocks that are only evaluated if needed (or in the case of the `switch` operator, an array of implicit Blocks for both operands).

Logical operators and methods in **dish** serve as the control structures for the language. **Note bene**: Be careful with logical chains such as `(condition).then(0).else(1)`--falsy values like 0 force the `else` block into execution regardless of the condition's truthiness. Therefore, to avoid any traps of falsy values, the ternary `[condition].switch[0, 1]` is preferred.

Similarly, to create traditional if/elseif/else logic in **dish**, use `[cond1, cond2].switch[0, 1, 2]`. The conditions Array is executed in sequence until it finds a truthy value (with all subsequent conditions being short-cictuited), and only the implicit Logic block from the right-hand Array corresponding to the same truthy index will be executed (or the last index if no truthy conditions were found).

| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `&&`     | Implicit Boolean `&&` Any   | `a.then(b)`       | `3 && 2`                   | `2`                   |
|          |                             |                   | `true && 0`                | `0`                   |
|          |                             |                   | `false && 2`               | `false`               |
|          |                             |                   |                            |                       |
| `\|\|`   | Implicit Boolean `\|\|` Any | `a.else(b)`       | `3 \|\| 2`                 | `3`                   |
|          |                             |                   | `false \|\| 0`             | `0`                   |
|          |                             |                   | `0 \|\| true`              | `true`                |
|          |                             |                   |                            |                       |
| `^^`     | Implicit Boolean `^^` Any   | `a.xor(b)`        | `3 ^^ 0`                   | `3`                   |
|          |                             |                   | `false ^^ 2`               | `2`                   |
|          |                             |                   | `3 ^^ 2`                   | `false`               |
|          |                             |                   |                            |                       |
| `?`      | Boolean Array `?` Any Array | `a.switch[b]`     | `[1 > 0] ? [3, 4]`         | `3`                   |
|          |                             |                   | `[1 < 0] ? [3, 4]`         | `4`                   |
|          |                             |                   |                            |                       |
| `??`     | Implicit Boolean `??` Any   | `a.redo(b)`       | `(a < 9) ?? (++a)`         | `9`                   |
|          |                             |                   | `(1 < 0) ?? (++a)`         | `false`               |

## Comparison Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `==`     | Scalar `==` Implicit Number | `a.equals(b)`     | `2 == '2.0'`               | `true`                |
|          | String `==` Implicit String | `a.equals(b)`     | `'2.0' == 2`               | `false`               |
|          | List `==` Implicit List     | `a.equals(b)`     | `[1, 2, 3] == [1, 2, 3]`   | `true`                |
|          |                             |                   |                            |                       |
| `!=`     | Scalar `!=` Implicit Number | `a.isnt(b)`       | `2 != '2.0'`               | `false`               |
|          | String `!=` Implicit String | `a.isnt(b)`       | `'2.0' != 2`               | `true`                |
|          | List `!=` Implicit List     | `a.isnt(b)`       | `[1, 2, 3] != [1, 2, 3]`   | `false`               |
|          |                             |                   |                            |                       |
| `>`      | Scalar `>` Implicit Number  | `a.above(b)`      | `3 > '10'`                 | `false`               |
|          | String `>` Implicit String  | `a.above(b)`      | `'3' > 10`                 | `true`                |
|          | List `>` Implicit Number    | `a.above(b)`      | `[1, 1, 1] > [2, 2]`       | `true`                |
|          |                             |                   |                            |                       |
| `<`      | Scalar `<` Implicit Number  | `a.below(b)`      | `3 < '10'`                 | `true`                |
|          | String `<` Implicit String  | `a.below(b)`      | `'3' < 10`                 | `false`               |
|          | List `<` Implicit Number    | `a.below(b)`      | `[1, 1, 1] < [2, 2]`       | `false`               |
|          |                             |                   |                            |                       |
| `>=`     | Scalar `>=` Implicit Number | `a.over(b)`       | `3 >= '10'`                | `false`               |
|          | String `>=` Implicit String | `a.over(b)`       | `'3' >= 10`                | `true`                |
|          | List `>=` Implicit Number   | `a.over(b)`       | `[1, 1] >= [2, 2]`         | `true`                |
|          |                             |                   |                            |                       |
| `<=`     | Scalar `<=` Implicit Number | `a.under(b)`      | `3 <= '10'`                | `true`                |
|          | String `<=` Implicit String | `a.under(b)`      | `'3' <= 10`                | `false`               |
|          | List `<=` Implicit Number   | `a.under(b)`      | `[1, 1] <= [2, 2]`         | `true`                |

## Assignment Operators
When initializing a variable for local scope, use `:`; when initializing a variable for global scope, use `=`. After initialization, `=` (and all variants) should be used to re-assign existing global _or_ local variable.

The member assignment operator (`:=`) is a special assignment operator that allows you to assign a value to a variable member (an array index or hash property) and return the parent variable.

| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `:`      | Variable `:` Any            | `a.define(b)`     | `a: 5`                     | `5`                   |
|          |                             |                   |                            |                       |
| `:`      | `:`Variable                 |                   | `:a(...)`                  | // Block argument     |
| (Unary)  |                             |                   |                            |                       |
|          |                             |                   |                            |                       |
| `=`      | Variable `=` Any            | `a.assign(b)`     | `a = [1, 2, 3]`            | `[1, 2, 3]`           |
| `:=`     | Member `:=` Any             | `a.0.replace(b)`  | `a = [1, 2, 3], a.0 := 4`  | `[4, 2, 3]`           |
| `~=`     |                             |                   | `a = 'bin', a ~= 'ary'`    | `"binary"`            |
| `+=`     |                             |                   | `a = 1, a += 1`            | `2`                   |
| `-=`     |                             |                   | `a = 1, a -= 1`            | `0`                   |
| `*=`     |                             |                   | `a = 2, a *= 2`            | `4`                   |
| `/=`     |                             |                   | `a = 6, a /= 2`            | `3`                   |
| `%=`     |                             |                   | `a = 5, a %= 4`            | `1`                   |
| `&=`     |                             |                   | `a = 5, a &= 3`            | `1`                   |
| `^=`     |                             |                   | `a = 5, a ^= 3`            | `6`                   |
| `\|=`    |                             |                   | `a = 5, a \|= 3`           | `7`                   |
|          |                             |                   |                            |                       |
| `++`     | `++`Scalar                  | `a.increment`     | `a = 8, ++a`               | `9`                   |
| (Unary)  | `++`Strin                   | `a.increment`     | `a = 'binary', ++a`        | `"cjobsz"`            |
|          | `++`Array                   | `a.increment`     | `a = [1, 2, 3], ++a`       | `[2, 3, 4]`           | 
|          | `++`Hash                    | `a.increment`     | `a = {x: 5}, ++a`          | `{"x": 6}`            |
|          |                             |                   |                            |                       |
| `--`     | `--`Scalar                  | `a.decrement`     | `a = 9, --a`               | `8`                   |
| (Unary)  | `--`String                  | `a.decrement`     | `a = 'cjobsz', --a`        | `"binary"`            |
|          | `--`Array                   | `a.decrement`     | `a = [1, 2, 3], --a`       | `[0, 1, 2]`           |
|          | `--`Hash                    | `a.decrement`     | `a = {x: 5}, --a`          | `{"x": 4}`            |
|          |                             |                   |                            |                       |
| `<<`     | Scalar [See Bitwise Ops](#bitwise-operators) |  |                            |                       |
|          | String [See Bitiwse Ops](#bitwise-operators) |  |                            |                       |
|          | Arary `<<` Implicit Array   | `a.push(b)`       | `a = [1, 2], a << 3`       | `[1, 2, 3]`           |
|          | Hash `<<` Implicit Hash     | `a.push(b)`       | `a = {x: 1}, a << {y: 2}`  | `{"x": 1, "y": 2}`    |
|          |                             |                   |                            |                       |
| `>>`     | Scalar [See Bitwise Ops](#bitwise-operators) |  |                            |                       |
|          | String [See Bitiwse Ops](#bitwise-operators) |  |                            |                       |
|          | Array `>>` Implicit Array   | `a.unshift(b)`    | `a = [1, 2], a >> 3`       | `[3, 1, 2]`           |
|          | Hash `>>` Implicit Hash     | `a.unshift(b)`    | `a = {x: 1}, a >> {y: 2}`  | `{"x": 1, "y": 2}`    |
|          |                             |                   |                            |                       |
| `<<`     | `<<`Scalar                  | `a.shift`         | `a = 12, [<<a, a]`         | `[0, 6]` i            |
| (Unary)  | `<<`String                  | `a.shift`         | `a = 'binary', [<<a, a]`   | `["b", "inary"]`      |
|          | `<<`List                    | `a.shift`         | `a = [1, 2, 3], [<<a, a]`  | `[1, [2, 3]]`         |
|          |                             |                   |                            |                       |
| `>>`     | `>>`Scalar                  | `a.pop`           | `a = 12, [>>a, a]`         | `[1, 4]`              |
| (Unary)  | `>>`String                  | `a.pop`           | `a = 'binary', [>>a, a]`   | `["y", "binar"]`      |
|          | `>>`List                    | `a.pop`           | `a = [1, 2, 3], [>>a, a]`  | `[3, [1, 2]]`         |


**Dish** also supports parallel assignment when the left-hand operand is an Array of variables, as in `[a, b, c] = [1, 2, 3]`. Parallel assigment is useful for swapping the values held in two variables without the need of a temporary holding variable: `[a, b] = [b, a]` is equivalent to `t = a, a = b, b = t`.

## Bitwise Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `&`      | Scalar `&` Implicit Number  | `a.intersect(b)`  | `5 & 3`                    | `1`                   |
|          | String `&` String           | `a.intersect(b)`  | `'codfish' & 'airshed'`    | `"dish"`              |
|          | String `&` Implicit Number  | `a.intersect(b)`  | `'5' & 3`                  | `1`                   |
|          | Array `&` Implicit Array    | `a.intersect(b)`  | `[1, 2, 3] & [4, 3, 2]`    | `[2, 3]`              |
|          | Hash `&` Implicit Hash      | `a.intersect(b)`  | `{x: 1} & {x: 2, y: 2}`    | `{"x": 2}`            |
|          |                             |                   |                            |                       |
| `^`      | Scalar `^` Implicit Number  | `a.exclude(b)`    | `5 ^ 3`                    | `6`                   |
|          | String `^` String           | `a.exclude(b)`    | `'danish' ^ 'nana'`        | `"dish"`              |
|          | String `^` Implicit Number  | `a.exclude(b)`    | `'5' ^ 3`                  | `6`                   |
|          | Array `^` Implicit Array    | `a.exclude(b)`    | `[1, 2, 3] ^ [4, 3, 2]`    | `[1, 4]`              |
|          | Hash `^` Implicit Hash      | `a.exclude(b)`    | `{x: 1} ^ {y: 2}`          | `{}`                  |
|          |                             |                   |                            |                       |
| `^`      | `^`Scalar                   | `a.invert`        | `^5`                       | `-6`                  |
| (Unary)  | `^`String                   | `a.invert`        | `^'binary'`                | `������`              |
|          | `^`List                     | `a.invert`        | `^[1, 2, 3]`               | `[-2, -3, -4]`        |
|          |                             |                   |                            |                       |
| `\|`     | Scalar `\|` Implicit Number | `a.union(b)`      | `5 \| 3`                   | `7`                   |
|          | String `\|` String          | `a.union(b)`      | `'didi' \| 'sihi'`         | `"dish"`              |
|          | String `\|` Implicit Number | `a.union(b)`      | `'5' \| 3`                 | `7`                   |
|          | Array `\|` Implicit Array   | `a.union(b)`      | `[1, 2, 3] \| [4, 2, 1]`   | `[1, 2, 3, 4]`        |
|          | Hash `\|` Implicit Hash     | `a.union(b)`      | `{x: 1} \| {x: 2}`         | `{"x": 2}`            |
|          |                             |                   |                            |                       |
| `<<`     | Scalar `<<` Implicit Number | `a.lshift(b)`     | `5 << 3`                   | `40`                  |
|          | String `<<` Implicit Number | `a.rpad(b)`       | `'binary' << 3`            | `"binary   "`         |
|          | Array [See Assignment Ops](#assignment-operators) |                        | |                       |
|          | Hash [See Assignment Ops](#assignment-operators)  |                        | |                       |
|          |                             |                   |                            |                       |
| `>>`     | Scalar `>>` Implicit Number | `a.rshift(b) `    | `40 >> 3`                  | `5`                   |
|          | String `>>` Implicit Number | `a.lpad(b)`       | `'binary' >> 3`            | `"   binary"`         |
|          | Array [See Assignment Ops](#assignment-operators) |                        | |                       |
|          | Hash [See Assignment Ops](#assignment-operators)  |                        | |                       |

## Binary Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `~`      | Scalar `+` Implicit String  | `a.concat(b)`     | `10 ~ 20`                  | `"1020"`              |
|          | String `+` Implicit String  | `a.concat(b)`     | `'bin' ~ 'ary'`            | `"binary"`            |
|          | List `+` Implicit String    | `a.join(b)`       | `[1, 2, 3] ~ ' '`          | `"1 2 3"`             |
|          |                             |                   |                            |                       |
| `+`      | Scalar `+` Implicit Number  | `a.add(b)`        | `10 + 20`                  | `30`                  |
|          | String `+` String           | `a.add(b)`        | `'bin' + 'ary'`            | `"binary"`            |
|          | String `+` Implicit Number  | `a.cipher(b)`     | `'BINARY' + 20`            | `"V]bUfm"`            |
|          | Array `+` Implicit Array    | `a.add(b)`        | `[1, 2] + [2, 3]`          | `[1, 2, 2, 3]`        |
|          | Hash `+` Implicit Hash      | `a.add(b)`        | `{x: 1} + {y: 2}`          | `{"x": 1, "y": 2}`    |
|          |                             |                   |                            |                       |
| `-`      | Scalar `-` Implicit Number  | `a.subtract(b)`   | `10 - 20`                  | `-10`                 |
|          | String `-` String           | `a.remove(b)`     | `'binary' - 'ary'`         | `"bin"`               |
|          | String `-` Implicit Number  | `a.decipher(b)`   | `'binary' - 20`            | `"NUZM^e"`            |
|          | Array `-` Implicit Array    | `a.remove(b)`     | `[1, 2, 2, 3, 4] - [2, 3]` | `[1, 2, 4]`           |
|          | Hash `-` Implicit Hash      | `a.remove(b)`     | `{x: 1, y: 2} - {"x"}`     | `{"y": 2}`            |
|          |                             |                   |                            |                       |
| `*`      | Scalar `*` :(Block)         | `a.times(b)`      | `4 * :(++a)`               | `[1, 2, 3, 4]`        |
|          | Scalar `*` Implicit Number  | `a.multiply(b)`   | `10 * 20`                  | `200`                 |
|          | String `*` :(Block)         | `a.map:c(...)`    | `'binary' * :c(c + 1)`     | `"cjobsz"`            |
|          | String `*` Implicit Number  | `a.repeat(b)`     | `'a' * 5`                  | `"aaaaa"`             |
|          | List `*` :(Block)           | `a.map:i(...)`    | `[1, 2, 3] * :i(i * 2)`    | `[1, 4, 6]`           |
|          | List `*` List               | `a.dot(b)`        | `[1, 2, 3] * [1, 2, 3]`    | `[1, 4, 9]`           |
|          | List `*` Implicit Number    | `a.repeat(b)`     | `[1, 2, 3] * 2`            | `[1, 2, 3, 1, 2, 3]`  |
|          |                             |                   |                            |                       |
| `/`      | Scalar `/` Implicit Number  | `a.divide(b)`     | `10 / 20`                  | `0.5`                 |
|          | String `/` :(Block)         | `a.split:i(...)`  | `'binary' / :i(i.vowel)`   | `["b", "n", "ry"]`    |
|          | String `/` String           | `a.split(b)`      | `'binary' / 'in'`          | `["b", "ary"]`        |
|          | String `/` Implicit Number  | `a.split(b)`      | `'binary' / 2`             | ` ["bi", "na", "ry"]` |
|          | List `/` :(Block)           | `a.split:i(...)`  | `[2, 2, 3, 4] / :i(i % 2)` | `[[2, 2], [4]]`       |
|          | List `/` Implicit Number    | `a.split(b)`      | `[2, 2, 3, 4] / 2`         | `[[2, 2], [3, 4]]`    |
|          |                             |                   |                            |                       |
| `%`      | Scalar `%` Implicit Number  | `a.remainder(b)`  | `20 % 7`                   | `6`                   |
|          | String `%` :(Block)         | `a.filter:i(...)` | `'binary' % :i(i.vowel)`   | `"ia"`                |
|          | String `%` Implicit Number  | `a.every(b)`      | `'binary' % 2`             | `"bnr"`               |
|          | List `%` :(Block)           | `a.filter:i(...)` | `[1, 1, 2, 3] % :i(i % 2)` | `[1, 1, 3]`           |
|          | List `%` Implicit Number    | `a.every(b)`      | `[1, 1, 2, 3] % 2`         | `[1, 2]`              |
|          |                             |                   |                            |                       |
| `**`     | Scalar `**` Implicit Number | `a.pow(b)`        | `2 ** 3`                   | `8`                   |
|          | String `**` :(Block)        | `a.sort:x:y(...)` | `'binary' ** :x:y(y < x)`  | `"yrniba"`            | 
|          | String `**` Implicit Number | `a.pow(b)`        | `'2' ** '3'`               | `8`                   |
|          | List `**` :(Block)          | `a.sort:x:y(...)` | `[7, 9, 4] ** :x:y(y < x)` | `[9, 7, 4]`           |
|          | List `**` Implicit Number   | `a.rotate(b)`     | `[7, 9, 4] ** 1`           | `[4, 7, 9]`           |
|          |                             |                   |                            |                       |
| `@`      | Scalar `@` Implicit Number  | `a.round(b)`      | `3.14 @ 1`                 | `3.1`                 |
|          | String `@` :(Block)         | `a.index:i(...)`  | `'binary' @ :i(i.vowel)`   | `[1, 3]`              |
|          | String `@` Implicit String  | `a.index(b)`      | `'binary' @ 'ary'`         | `[3]`                 |
|          | Array `@` :(Block)          | `a.index:i(...)`  | `[1, 2, 3] @ :i(i % 2)`    | `[0, 2]`              |
|          | Array `@` Any               | `a.index(b)`      | `[1, 2, 3] @ 3`            | `[2]`                 |
|          | Hash `@` :(Block)           | `a.index:i(...)`  | `{f: 1, b: 2} @ :i(i % 2)` | `["f"]`               |
|          | Hash `@` Any                | `a.index(b)`      | `{a: 1, b: 2, c: 2} @ 2`   | `["b", "c"]`          |
|          |                             |                   |                            |                       |
| `~~`     | Scalar `~~` String          | `a.fmt(b)`        | `9 ~~ '%.2f'`              | `"9.00"`              |
|          | Scalar `~~` Implicit Number | `a.base(b)`       | `9 ~~ 2`                   | `"1001"`              |
|          | String `~~` String          | `a.fmt(b)`        | `'binary' ~~ '%7s'`        | `" binary"`           |
|          | String `~~` Implicit Number | `a.base(b)`       | `'1001' ~~ 2`              | `9`                   |
|          | Array `~~` String           | `a.fmt(b)`        | `[1, 2] ~~ '%.1f'`         | `["1.0", "2.0"]`      |
|          | Array `~~` Implicit Number  | `a.base(b)`       | `[1, 2, 3] ~~ 2`           | `["1", "10", "11"]`   |
|          | Hash `~~` String            | `a.fmt(b)`        | `{x: 1} ~~ '%.1f'`         | `{"x": "1.0"}`        |
|          | Hash `~~` Implicit Number   | `a.base(b)`       | `{a: 3} ~~ 2`              | `{"a": "11"}`         |
|          |                             |                   |                            |                       |
| `..`     | Scalar `..` String          | `a.to(b)`         | `95 .. 'a'`                | `[95, 96, 97]`        |
|          | Scalar `..` Implicit Number | `a.to(b)`         | `3 .. 1`                   | `[3, 2, 1]`           |
|          | String `..` String          | `a.to(b)`         | `'a' .. 'c'`               | `["a", "b", "c"]`     |
|          | String `..` Implicit Number | `a.to(b)`         | `'a' .. 99`                | `["a", "b", "c"]`     |

**Note bene**: the `..` range/to operator supports positve as well as negative ranges (as forced integers for numbers). In cases where a range is used short-circuit a traditional for loop, you may need to first check that the right-hand operand is greater than the left-hand operand. Compare the following:

`dish -e 'a = 1; (0..3).map:i((1..i).map:j(a *= j)); a'` outputs `0`

`dish -e 'a = 1; (0..3).map:i((i > 1).then((1..i).map:j(a *= j))); a'` outputs `12`

## Unary Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `!`      | `!`Implicit Boolean         | `a.not`           | `!5`                       | `false`               |
|          |                             |                   |                            |                       |
| `?`      | `?`Implicit Boolean         | `a.bool`          | `?5`                       | `true`                |
|          |                             |                   |                            |                       |
| `-`      | `-`Scalar                   | `a.negate`        | `-5`                       | `-5`                  |
|          | `-`String                   | `a.negate`        | `-"5"`                     | `-5`                  |
|          | `-`List                     | `a.negsum`        | `-[1, 2, 3]`               | `-6`                  |
|          |                             |                   |                            |                       |
| `+`      | `+`Scalar                   | `a.num`           | `+5`                       | `5`                   |
|          | `+`String                   | `a.num`           | `+"5"                      | `5`                   |
|          | `+`List                     | `a.sum`           | `+[1, 2, 3]`               | `6`                   |
|          |                             |                   |                            |                       |
| `~`      | `~`Scalar                   | `a.str`           | `~5`                       | `"5"`                 |
|          | `~`String                   | `a.str`           | `~"5"`                     | `"5"`                 |
|          | `~`List                     | `a.join`          | `~[1, 2, 3]`               | `"123"`               |
|          |                             |                   |                            |                       |
| `/`      | `~`Scalar                   | `a.ratio`         | `/3.14`                    | `[157, 50]`           |
|          | `~`String                   | `a.arr`           | `/'bin'`                   | `["b", "i", "n"]`     |
|          | `~`List                     | `a.vals`          | `/{x: 1, y: 2}`            | `[1, 2]`              |
|          |                             |                   |                            |                       |
| `%`      | `%`Scalar                   | `a.hash`          | `%5`                       | `{"5": 5}`            |
|          | `%`String                   | `a.hash`          | `%'x'`                     | `{"x": "x"}`          |
|          | `%`List                     | `a.hash`          | `%[1, 2]`                  | `{"1": 1, "2": 2}`    |
|          |                             |                   |                            |                       |
| `@`      | `@`Scalar                   | `a.reverse`       | `@3.14`                    | `41.3`                |
|          | `@`String                   | `a.reverse`       | `@'binary'`                | `"yranib"`            |
|          | `@`Array                    | `a.reverse`       | `@[1, 2, 3]`               | `[3, 2, 1]`           |
|          | `@`Hash                     | `a.keys`          | `@{x: 1, y: 2}`            | `["x", "y"]`          |
|          |                             |                   |                            |                       |
| `#`      | `#`Scalar                   | `a.len`           | `#15`                      | `null`                |
|          | `#`String                   | `a.len`           | `#'binary'`                | `6`                   |
|          | `#`List                     | `a.len`           | `#[1, 2, 3]`               | `3`                   |
|          |                             |                   |                            |                       |
| `++`     | `++`Scalar                  | `a.increment`     | `++15`                     | `16`                  |
|          | `++`String                  | `a.increment`     | `++'15'`                   | `16`                  |
|          | `++`List                    | `a.increment`     | `++[1, 2, 3]`              | `[2, 3, 4]`           |
|          | [Also See Assignment Ops](#assignment-operators) |                         | |                       |
|          |                             |                   |                            |                       |
| `--`     | `--`Scalar                  | `a.decrement`     | `--15`                     | `14`                  |
|          | `--`String                  | `a.decrement`     | `--'15'`                   | `14`                  |
|          | `--`List                    | `a.decrement`     | `--[1, 2, 3]`              | `[0, 1, 2]`           |
|          | [Also See Assignment Ops](#assignment-operators) |                         | |                       |
|          |                             |                   |                            |                       |
| `<<`     | `<<`Scalar                  | `a.first`         | `<<15`                     | `15`                  |
|          | `<<`String                  | `a.first`         | `<<'binary'`               | `"b"`                 |
|          | `<<`List                    | `a.first`         | `<<[1, 2, 3]`              | `1`                   |
|          | [Also See Assignment Ops](#assignment-operators) |                         | |                       |
|          |                             |                   |                            |                       |
| `>>`     | `>>`Scalar                  | `a.last`          | `>>15`                     | `15`                  |
|          | `>>`String                  | `a.last`          | `>>'binary'`               | `"y"`                 |
|          | `>>`List                    | `a.last`          | `>>[1, 2, 3]`              | `3`                   |
|          | [Also See Assignment Ops](#assignment-operators) |                         | |                       |
|          |                             |                   |                            |                       |
| `*`      | `*`Scalar                   | `a.prime`         | `*7`                       | `true`                |
|          | `*`String                   | `a.eval`          | `*'2 + 2'`                 | `4`                   |
|          | `*`List                     | `a.product`       | `*[1, 2, 3, 4]`            | `24`                  |
|          |                             |                   |                            |                       |
| `**`     | `**`Scalar                  | `a.divisors`      | `**8`                      | `[1, 2, 4]`           |
|          | `**`String                  | `a.sort`          | `**'binary'`               | `"abinry"`            |
|          | `**`List                    | `a.sort`          | `**[10, 1, 2, 4]`          | `[1, 2, 4, 10]`       |
|          |                             |                   |                            |                       |
| `<`      | `<`Scalar                   | `a.floor`         | `<3.14`                    | `3`                   |
|          | `<`String                   | `a.lc`            | `<'BINARY'`                | `"binary"`            |
|          | `<`List                     | `a.min`           | `<[1, 2, 3]`               | `1`                   |
|          |                             |                   |                            |                       |
| `>`      | `>`Scalar                   | `a.ceil`          | `>3.14`                    | `4`                   |
|          | `>`String                   | `a.uc`            | `>'binary'`                | `"BINARY"`            |
|          | `>`List                     | `a.max`           | `>[1, 2, 3]`               | `3`                   |
|          |                             |                   |                            |                       |
| `=`      | `=`Scalar                   | `a.int`           | `=3.14`                    | `3`                   |
|          | `=`String                   | `a.avg`           | `='binary'`                | `"k"`                 |
|          | `=`List                     | `a.avg`           | `=[1, 2, 3]`               | `2`                   |
|          |                             |                   |                            |                       |
| `~~`     | `~~`Scalar                  | `a.chr`           | `~~65`                     | `"A"`                 |
|          | `~~`String                  | `a.ord`           | `~~'f'`                    | `102`                 |
|          | `~~`List                    | `a.ascii`         | `~~[65, 'B', 67]`          | `["A", 65, "C"]`      |
|          |                             |                   |                            |                       |
| `^`      | [See Bitwise Ops](#bitwise-operators)           | |                          |                       |
|          |                             |                   |                            |                       |
| `&`      | `&`Scalar                   | `a.defined`       | `&a`                       | `false`               |
|          | `&`String                   | `a.defined`       | `a = 'x', &a`              | `true`                | 
|          | `&`Array                    | `a.compact`       | `&[0, 1, 2]`               | `[1, 2]`              |
|          | `&`Hash                     | `a.compact`       | `&{a: 0, b: 1, c: null}`   | `{"b": 1}`            |
|          |                             |                   |                            |                       |
| `\|`     | `\|`Scalar                  | `a.abs`           | `\|-3`                     | `3`                   |
|          | `\|`String                  | `a.uniq`          | `\|'bookkeeper'`           | `"bokepr"`            |
|          | `\|`List                    | `a.uniq`          | `\|[1, 1, 2, 3, 3]`        | `[1, 2, 3]`           |

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
| `&&`     | Implicit Boolean `&&` Any   | `a.then(b)`       | `3 && 2`                   | `2`                   |
|          |                             |                   | `true && 0`                | `false`               |
|          |                             |                   |                            |                       |
| `\|\|`   | Implicit Boolean `\|\|` Any | `a.else(b)`       | `3 \|\| 2`                 | `3`                   |
|          |                             |                   | `false \|\| 0`             | `0`                   |
|          |                             |                   |                            |                       |
| `^^`     | Implicit Boolean `^^` Any   | `a.xor(b)`        | `3 ^^ 2`                   | `false`               |
|          |                             |                   | `false ^^ 2`               | `2`                   |
|          |                             |                   |                            |                       |
| `?`      | Boolean Array `?` Any Array | `a.switch[b]`     | `[1 > 0] ? [3, 4]`         | `3`                   |
|          |                             |                   | `[1 < 0] ? [3, 4]`         | `4`                   |
|          |                             |                   |                            |                       |
| `??`     | Implicit Boolean `??` Any   | `a.redo(b)`       | `(a < 9) ?? (++a)`         | `9`                   |
|          |                             |                   | `(1 < 0) ?? (++a)`         | `false`               |

## Comparison Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `==`     | Scalar `==` Implicit Number | `a.equals(b)`     | `2 == 2`                   | `true`                |
|          | String `==` Implicit String | `a.equals(b)`     | `'2' == '2.0'`             | `false`               |
|          | List `==` Implicit List     | `a.equals(b)`     | `[1, 2, 3] == [1, 2, 3]`   | `true`                |
|          |                             |                   |                            |                       |
| `!=`     | Scalar `!=` Implicit Number | `a.isnt(b)`       | `2 != 2`                   | `false`               |
|          | String `!=` Implicit String | `a.isnt(b)`       | `'2' != '2.0'`             | `true`                |
|          | List `!=` Implicit List     | `a.isnt(b)`       | `[1, 2, 3] != [1, 2, 3]`   | `false`               |
|          |                             |                   |                            |                       |
| `>`      | Scalar `>` Implicit Number  | `a.above(b)`      | `3 > 10`                   | `false`               |
|          | String `>` String           | `a.above(b)`      | `'3' > '10'`               | `true`                |
|          | String `>` Implicit Number  | `a.above(b)`      | `'3' > 10`                 | `false`               |
|          | List `>` Implicit Number    | `a.above(b)`      | `[1, 1, 1] > [2, 2]`       | `true`                |
|          |                             |                   |                            |                       |
| `<`      | Scalar `<` Implicit Number  | `a.below(b)`      | `3 < 10`                   | `true`                |
|          | String `<` String           | `a.below(b)`      | `'3' < '10'`               | `false`               |
|          | String `<` Implicit Number  | `a.below(b)`      | `'3' < 10`                 | `true`                |
|          | List `<` Implicit Number    | `a.below(b)`      | `[1, 1, 1] < [2, 2]`       | `false`               |
|          |                             |                   |                            |                       |
| `>=`     | Scalar `>=` Implicit Number | `a.over(b)`       | `3 >= 10`                  | `false`               |
|          | String `>=` String          | `a.over(b)`       | `'3' >= '10'`              | `true`                |
|          | String `>=` Implicit Number | `a.over(b)`       | `'3' >= 10`                | `false`               |
|          | List `>=` Implicit Number   | `a.over(b)`       | `[1, 1] >= [2, 2]`         | `true`                |
|          |                             |                   |                            |                       |
| `<=`     | Scalar `<=` Implicit Number | `a.under(b)`      | `3 <= 10`                  | `true`                |
|          | String `<=` String          | `a.under(b)`      | `"3" <= "10"`              | `false`               |
|          | String `<=` Implicit Number | `a.under(b)`      | `'3' <= 10`                | `true`                |
|          | List `<=` Implicit Number   | `a.under(b)`      | `[1, 1] <= [2, 2]`         | `true`                |

## Assignment Operators
When initializing a variable for local scope, use `:`; when initializing a variable for global scope, use `=`. After initialization, `=` (and all variants) should be used to re-assign existing global _or_ local variable.

The member assignment operator (`:=`) is a special assignment operator that allows you to assign a value to a variable member (an array index or hash property) and return the parent variable.

| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `:`      | Variable `:` Any            | `a.define(b)`     | `a: 5`                     | `5`                   |
|          |                             |                   |                            |                       |
| `=`      | Variable `=` Any            | `a.assign(b)`     | `a = [1, 2, 3]`            | `[1, 2, 3]`           |
| `:=`     | Member `:=` Any             |                   | `a = [1, 2, 3], a.0 := 4`  | `[4, 2, 3]`           |
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
|          | `++`Hash                    | `a.increment`     | `a = {foo: 5}, ++a`        | `{"foo": 6}`          |
|          |                             |                   |                            |                       |
| `--`     | `--`Scalar                  | `a.decrement`     | `a = 9, --a`               | `8`                   |
| (Unary)  | `--`String                  | `a.decrement`     | `a = 'cjobsz', --a`        | `"binary"`            |
|          | `--`Array                   | `a.decrement`     | `a = [1, 2, 3], --a`       | `[0, 1, 2]`           |
|          | `--`Hash                    | `a.decrement`     | `a = {foo: 5}, --a`        | `{"foo": 4}`          |
|          |                             |                   |                            |                       |
| `<<`     | Scalar [See Bitwise Ops](#bitwise-operators) |  |                            |                       |
|          | String [See Bitiwse Ops](#bitwise-operators) |  |                            |                       |
|          | Arary `<<` Implicit Array   | `a.push(b)`       | `a = [1, 2], a << 3`       | `[1, 2, 3]`           |
|          | Hash `<<` Implicit Hash     | `a.push(b)`       | `a = {foo: 1}, a << 2`     | `{"foo": 1, "2": 2}`  |
|          |                             |                   |                            |                       |
| `>>`     | Scalar [See Bitwise Ops](#bitwise-operators) |  |                            |                       |
|          | String [See Bitiwse Ops](#bitwise-operators) |  |                            |                       |
|          | Array `>>` Implicit Array   | `a.unshift(b)`    | `a = [1, 2], a >> 3`       | `[3, 1, 2]`           |
|          | Hash `>>` Implicit Hash     | `a.unshift(b)`    | `a = {foo: 1}, a >> 2`     | `{"2": 2, "foo": 1}`  |
|          |                             |                   |                            |                       |
| `<<`     | `<<`Scalar                  | `a.shift`         | `a = 3, [<<a, a]`          | `[3, 0]`              |
| (Unary)  | `<<`String                  | `a.shift`         | `a = 'binary', [<<a, a]`   | `["b", "inary"]`      |
|          | `<<`List                    | `a.shift`         | `a = [1, 2, 3], [<<a, a]`  | `[1, [2, 3]]`         |
|          |                             |                   |                            |                       |
| `>>`     | `>>`Scalar                  | `a.pop`           | `a = 3, [>>a, a]`          | `[3, 0]`              |
| (Unary)  | `>>`String                  | `a.pop`           | `a = 'binary', [>>a, a]`   | `["y", "binar"]`      |
|          | `>>`List                    | `a.pop`           | `a = [1, 2, 3], [>>a, a]`  | `[3, [1, 2]]`         |


**Dish** also supports parallel assignment when the left-hand operand is an Array of variables, as in `[a, b, c] = [1, 2, 3]`. Parallel assigment is useful for swapping the values held in two variables without the need of a temporary holding variable: `[a, b] = [b, a]` is equivalent to `t = a, a = b, b = t`.

## Bitwise Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `&`      | Scalar `&` Implicit Number  | `a.intersect(b)`  | `5 & 3`                    | `1`                   |
|          | String `&` String           | `a.intersect(b)`  | `'jadh!' & '____'`         | `"JADH"`              |
|          | String `&` Implicit Number  | `a.intersect(b)`  | `'5' & 3`                  | `1`                   |
|          | Array `&` Implicit Array    | `a.intersect(b)`  | `[1, 2, 3] & [4, 3, 2]`    | `[2, 3]`              |
|          | Hash `&` Implicit Hash      | `a.intersect(b)`  | `{f: 1} & {f: 2, b: 2}`    | `{"f": 2}`            |
|          |                             |                   |                            |                       |
| `^`      | Scalar `^` Implicit Number  | `a.exclude(b)`    | `5 ^ 3`                    | `6`                   |
|          | String `^` String           | `a.exclude(b)`    | `'d O!' ^ ' I<I!'`         | `"Dish!"`             |
|          | String `^` Implicit Number  | `a.exclude(b)`    | `'5' ^ 3`                  | `6`                   |
|          | Array `^` Implicit Array    | `a.exclude(b)`    | `[1, 2, 3] ^ [4, 3, 2]`    | `[1, 4]`              |
|          | Hash `^` Implicit Hash      | `a.exclude(b)`    | `{foo: 1} ^ {foo: 2}`      | `{}`                  |
|          |                             |                   |                            |                       |
| `^`      | `^`Scalar                   | `a.invert`        | `^5`                       | `-6`                  |
| (Unary)  | `^`String                   | `a.invert`        | `^'binary'`                | `������`              |
|          | `^`List                     | `a.invert`        | `^[1, 2, 3]`               | `[-2, -3, -4]`        |
|          |                             |                   |                            |                       |
| `\|`     | Scalar `\|` Implicit Number | `a.union(b)`      | `5 \| 3`                   | `7`                   |
|          | String `\|` String          | `a.union(b)`      | `'JA' \| '  dh!'`          | `"jadh!"`             |
|          | String `\|` Implicit Number | `a.union(b)`      | `'5' \| 3`                 | `7`                   |
|          | Array `\|` Implicit Array   | `a.union(b)`      | `[1, 2, 3] \| [4, 2, 1]`   | `[1, 2, 3, 4]`        |
|          | Hash `\|` Implicit Hash     | `a.union(b)`      | `{foo: 1} \| {foo: 2}`     | `{"foo": 2}`          |
|          |                             |                   |                            |                       |
| `<<`     | Scalar `<<` Implicit Number | `a.lshift(b)`     | `5 << 3`                   | `40`                  |
|          | String `<<` Implicit Number | `a.lshift(b)`     | `'5' << 3`                 | `40`                  |
|          | Array [See Assignment Ops](#assignment-operators) |                        | |                       |
|          | Hash [See Assignment Ops](#assignment-operators)  |                        | |                       |
|          |                             |                   |                            |                       |
| `>>`     | Scalar `>>` Implicit Number | `a.rshift(b) `    | `40 >> 3`                  | `5`                   |
|          | String `>>` Implicit Number | `a.rshift(b)`     | `'40' >> 3`                | `5`                   |
|          | Array [See Assignment Ops](#assignment-operators) |                        | |                       |
|          | Hash [See Assignment Ops](#assignment-operators)  |                        | |                       |

## Binary Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `~`      | Scalar `+` Implicit String  | `a.concat(b)`     | `10 ~ 20`                  | `"1020"`              |
|          | String `+` Implicit String  | `a.concat(b)`     | `'10' ~ '20'`              | `"1020"`              |
|          | List `+` Implicit String    | `a.join(b)`       | `[1, 2, 3] ~ ' '`          | `"1 2 3"`             |
|          |                             |                   |                            |                       |
| `+`      | Scalar `+` Implicit Number  | `a.add(b)`        | `10 + 20`                  | `30`                  |
|          | String `+` Implicit Number  | `a.add(b)`        | `'10' + '20'`              | `30`                  |
|          | Array `+` Implicit Array    | `a.add(b)`        | `[1, 2] + [2, 3]`          | `[1, 2, 2, 3]`        |
|          | Hash `+` Implicit Hash      | `a.add(b)`        | `{x: 1} + {y: 2}`          | `{"x": 1, "y": 2}`    |
|          |                             |                   |                            |                       |
| `-`      | Scalar `-` Implicit Number  | `a.subtract(b)`   | `10 - 20`                  | `-10`                 |
|          | String `-` Implicit Number  | `a.subtract(b)`   | `'10' - '20'`              | `-10`                 |
|          | Array `-` Implicit Array    | `a.remove(b)`     | `[1, 2] - [2, 3]`          | `[1]`                 |
|          | Hash `-` Implicit Hash      | `a.remove(b)`     | `{x: 1, y: 2} - {"x"}`     | `{"y": 2}`            |
|          |                             |                   |                            |                       |
| `*`      | Scalar `*` Implicit Number  | `a.multiply(b)`   | `10 * 20`                  | `200`                 |
|          | String `*` Implicit Number  | `a.repeat(b)`     | `'a' * 5`                  | `"aaaaa"`             |
|          | List `*` :(Block)           | `a.map:i(...)`    | `[1, 2, 3] * :i(i + 1)`    | `[2, 3, 4]`           |
|          | List `*` List               | `a.dot(b)`        | `[1, 2] * [3, 4]`          | `[3, 8]`              |
|          | List `*` Implicit Number    | `a.repeat(b)`     | `[2] * 5`                  | `[2, 2, 2, 2, 2]`     |
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
|          | String `@` Implicit String  | `a.index(b)`      | `'binary' @ 'ary'          | `[3]`                 |
|          | Array `@` :(Block)          | `a.index:i(...)`  | `[1, 2, 3] @ :i(i % 2)`    | `[0, 2]`              |
|          | Array `@` Any               | `a.index(b)`      | `[1, 2, 3] @ 3             | `[2]`                 |
|          | Hash `@` :(Block)           | `a.index:i(...)`  | `{f: 1, b: 2} @ :i(i % 2)` | `["f"]`               |
|          | Hash `@` Any                | `a.index(b)`      | `{a: 1, b: 2, c: 2} @ 2`   | `["b", "c"]`          |
|          |                             |                   |                            |                       |
| `~~`     | Scalar `~~` String          | `a.fmt(b)`        | `9 ~~ '%.02f'`             | `"9.00"`              |
|          | Scalar `~~` Implicit Number | `a.base(b)`       | `9 ~~ 2`                   | `"1001"`              |
|          | String `~~` String          | `a.fmt(b)`        | `'binary' ~~ '%10s'`       | `"    binary"`        |
|          | String `~~` Implicit Number | `a.base(b)`       | `'1001' ~~ 2`              | `9`                   |
|          | Array `~~` String           | `a.fmt(b)`        | `[1, 2] ~~ '%.01f'`        | `["1.0", "2.0"]`      |
|          | Array `~~` Implicit Number  | `a.base(b)`       | `[1, 2, 3] ~~ 2            | `["1", "10", "11"]`   |
|          | Hash `~~` String            | `a.fmt(b)`        | `{foo: 1} ~~ '%.01f'`      | `{"foo": "1.0"}`      |
|          | Hash `~~` Implicit Number   | `a.base(b)`       | `{a: 3} ~~ 2`              | `{"a": "11"}`         |
|          |                             |                   |                            |                       |
| `..`     | Scalar `..` String          | `a.to(b)`         | `95 .. 'a'`                | `[95, 96, 97]`        |
|          | Scalar `..` Implicit Number | `a.to(b)`         | `3 .. 1`                   | `[3, 2, 1]`           |
|          | String `..` String          | `a.to(b)`         | `'a' .. 'c'`               | `["a", "b", "c"]`     |
|          | String `..` Implicit Number | `a.to(b)`         | `'a' .. 99`                | `["a", "b", "c"]`     |

## Unary Operators
| Operator | Operands                    | Method Name       | Example                    | Result                |
| -------- | --------------------------- | ----------------- | -------------------------- | --------------------- |
| `!`      | `!`Implicit Boolean         | `a.not`           | `!5`                       | `false`               |
|          |                             |                   |                            |                       |
| `?`      | `?`Implicit Boolean         | `a.bool`          | `?5`                       | `true`                |
|          |                             |                   |                            |                       |
| `-`      | `-`Scalar                   | `a.negate`        | `-5`                       | `-5`                  |
|          | `-`String                   | `a.negate`        | `-"5"`                     | `-5`                  |
|          | `-`List                     | `a.negate`        | `-[1, 2, 3]`               | `-6`                  |
|          |                             |                   |                            |                       |
| `+`      | `+`Scalar                   | `a.num`           | `+5`                       | `5`                   |
|          | `+`String                   | `a.num`           | `+"5"                      | `5`                   |
|          | `+`List                     | `a.sum`           | `+[1, 2, 3]`               | `6`                   |
|          |                             |                   |                            |                       |
| `~`      | `~`Scalar                   | `a.str`           | `~5`                       | `"5"`                 |
|          | `~`String                   | `a.str`           | `~"5"`                     | `"5"`                 |
|          | `~`List                     | `a.str`           | `~[1, 2, 3]`               | `"123"`               |
|          |                             |                   |                            |                       |
| `/`      | `~`Scalar                   | `a.ratio`         | `/3.14`                    | `[157, 50]`           |
|          | `~`String                   | `a.arr`           | `/'foo'`                   | `["f", "o", "o"]`     |
|          | `~`List                     | `a.arr`           | `/{foo: 1, bar: 2}`        | `[1, 2]`              |
|          |                             |                   |                            |                       |
| `%`      | `%`Scalar                   | `a.hash`          | `%5`                       | `{"5": 5}`            |
|          | `%`String                   | `a.hash`          | `%'foo'`                   | `{"foo": "foo"}`      |
|          | `%`List                     | `a.hash`          | `%[1, 2]`                  | `{"1": 1, "2": 2}`    |
|          |                             |                   |                            |                       |
| `@`      | `@`Scalar                   | `a.reverse`       | `@3.14`                    | `41.3`                |
|          | `@`String                   | `a.reverse`       | `@'binary'`                | `"yranib"`            |
|          | `@`Array                    | `a.reverse`       | `@[1, 2, 3]`               | `[3, 2, 1]`           |
|          | `@`Hash                     | `a.keys`          | `@{foo: 1, bar: 2}`        | `["foo", "bar"]`      |
|          |                             |                   |                            |                       |
| `#`      | `#`Scalar                   | `a.size`          | `#15`                      | `2`                   |
|          | `#`String                   | `a.length`        | `#'binary'`                | `6`                   |
|          | `#`List                     | `a.length`        | `#[1, 2, 3]`               | `3`                   |
|          |                             |                   |                            |                       |
| `^`      | [See Bitwise Ops](#bitwise-operators)           | |                          |                       |
|          |                             |                   |                            |                       |
| `++`     | `++`Scalar                  | `a.increment`     | `++15`                     | `16`                  |
|          | `++`String                  | `a.increment`     | `++'binary'`               | `"cjobsz"`            |
|          | `++`List                    | `a.increment`     | `++[1, 2, 3]`              | `[2, 3, 4]`           |
|          | [Also See Assignment Ops](#assignment-operators) |                         | |                       |
|          |                             |                   |                            |                       |
| `--`     | `--`Scalar                  | `a.decrement`     | `--15`                     | `14`                  |
|          | `--`String                  | `a.decrement`     | `--'cjobsz'`               | `"binary"`            |
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
|          | `*`String                   | `a.vowel`         | `*'a'`                     | `true`                |
|          | `*`List                     | `a.product`       | `*[1, 2, 3, 4]`            | `24`                  |
|          |                             |                   |                            |                       |
| `**`     | `**`Scalar                  | `a.divisors`      | `**8`                      | `[1, 2, 4]`           |
|          | `**`String                  | `a.sort`          | `**'binary'`               | `"abinry"`            |
|          | `**`List                    | `a.sort`          | `**[10, 1, 2, 4]`          | `[1, 2, 4, 10]`       |

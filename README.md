# dish
Dish is a dynamically-typed, keywordless programming language designed with the following goals:

* Allow quick and easy processing of input data (partiularly when the input is json data) with traditional programming operators (as opposed to the query-based operators of [jq](https://github.com/stedolan/jq#readme))

* Provide flexibility in language syntax so that developers can tackle tasks based on individual conceptualization (while still allowing peers to easily follow logical flow--that is, while avoiding the "write-only" criticsms of [perl](https://github.com/Perl/perl5#readme))

* Isolate behavioral functionality by treating all data types as objects (and remove control flow keywords, like `if`, `for`, and `while` in favor of object methods, which goes a step further than [ruby](https://github.com/ruby/ruby#readme))

* Allow for contextual parsing of variables and statements without requiring end-of-statement identifiers (while still avoiding the "where does this block end?" criticisms of [python](https://github.com/python/cpython#readme))

## Installation
Dish is interpreted by [Go](https://github.com/golang/go#readme). With Go installed, build the **dish** interpreter via `go build -o /usr/local/bin/dish` from the project root. You can then run **dish** files via `dish /path/to/file.dish` or with the `-e` command flag, as in `dish -e '["Hello World!"]'`

## Output
By default **dish** outputs the last evaluated statement to STDOUT. If the last evaluated statement is an [Array](#array) or a [Hash](#hash), the output is formatted as valid JSON. If the last evaluated statement is a Scalar, the scalar's raw output is printed to STDOUT.

## Syntax
Dish has 3 types of syntax blocks, Value Blocks represented by `()`, Array blocks represented by `[]`, and Hash blocks represented by `{}`. Expressions and statements in each block are terminated either by a statement ending newline<sup>*</sup> or by a comma (unless the newline or comma is encapsulated in a string literal).

Value blocks `()` return the last expression or statement in the block. A full **dish** program is inside an implicit value block when the first and last characters of the program are not `(` and `)` respectively. Array blocks `[]` and Hash blocks `{}` return the entire array or hash.

Comments in **dish** start with a non string-ecapsulated `#` character and end with a newline. There are no multiline comments in **dish**.

<sub>*Statement ending newlines are those not preceded by an opening block character or by a binary operator.</sub>

## Data Types and "Operators"
Dish has 5 data types (Booleans, Numbers, Strings, Arrays, and Hashes), all of which are objects with their own object methods. Traditional, symbolic operators (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods.

As such, all binary operators apply the right-hand operand as a parameter to a method on the left-hand object. In the expression `5 + 2`, `2` is a parameter of the `add` method on the Number object `5`, e.g. `5.add(2)`.

Similarly, traditional unary operators (which must be used as prefix operators when represented in shorthand notation<sup>*</sup>) correspond to methods with no parameters that are invoked on the single operand as an object. In the expression `!ready`, `ready` is the Boolean object on which the `not` method is invoked, e.g. `ready.not()` 

<sub>*Postfix `++` and `--` are not legal in **dish**</sub>

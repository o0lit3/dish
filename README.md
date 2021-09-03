# dish
**dish** is a dynamically-typed programming language designed with the following goals:

* Allow quick and easy processing of input data (partiularly when the input is JSON data) with traditional programming operators (as opposed to the query-based operators of [jq](https://github.com/stedolan/jq#readme))

* Provide flexibility in language syntax so that developers can tackle tasks based on individual conceptualization (while still allowing peers to easily follow logical flow--that is, while avoiding the "write-only" criticsms of [perl](https://github.com/Perl/perl5#readme))

* Isolate behavioral functionality by treating all data types as objects (and remove control flow keywords, like `if`, `for`, and `while` in favor of object methods, which goes a step further than [ruby](https://github.com/ruby/ruby#readme))

* Allow for contextual parsing of variables and statements without requiring end-of-statement identifiers (while avoiding the "where does this block end?" criticisms of [python](https://github.com/python/cpython#readme))

## Installation
**dish** is interpreted by [Go](https://github.com/golang/go#readme). With Go installed, build the **dish** interpreter via `go build -o /usr/local/bin/dish` from the project root. You can then run **dish** files via `dish /path/to/file.dish` or with the `-e` command flag, as in `dish -e '"Hello World!"'`

## Input
By default, **dish** places STDIN into a variable called `stdin`. If the data from STDIN is JSON, `stdin` is a data type representing that JSON data ([see Data Types](#data-types-and-operators)), otherwise `stdin` is an Array of STDIN lines.

As an example, the following curl/**dish** command will output a list of the last 5 github commits to the **dish** codebase:

```curl -s "https://api.github.com/repos/o0lit3/dish/commits?per_page=5" | dish -e 'stdin.map:data("$(data.commit.message):\n$(data.parents.0.html_url)").join("\n\n")'```

...and the following curl/**dish** command will output a list of unique github contributors in the last 100 commits of the **react** codebase:

```curl -s "https://api.github.com/repos/facebook/react/commits?per_page=100" | dish -e 'stdin.map:data(data.commit.author.name).uniq.join("\n")'```

## Output
By default, **dish** outputs the last evaluated statement to STDOUT. If the last evaluated statement is an Array or a Hash, the output is formatted as valid JSON. If the last evaluated statement is a Scalar, the scalar's raw output is printed to STDOUT.

This behavior allows you to pipe the output of one dish executable into another and use it as JSON input via `stdin` downstream:

`dish -e "[1, 2, 3]" | dish -e "stdin.sum"`

## [Data Types and Operators](tests/README.md)
For a detailed decription of operators, precedence, and implicit operator context, [read the operator documentation at tests/README.md](tests/README.md).

Like JSON itself, **dish** has 6 data types (Hash, Array, String, Number, Boolean, and Null), all of which are objects with their own object methods. Traditional, symbolic operators (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods.

As such, all binary operators apply the right-hand operand as a parameter to a method on the left-hand object. In the expression `5 + 2`, `2` is a parameter of the `add` method on the Number object `5`, e.g. `5.add(2)`.

Similarly, traditional unary operators (which are always prefix operators when represented in shorthand notation<sup>*</sup>) correspond to methods with no parameters that are invoked on the single operand as an object. In the expression `!ready`, `ready` is the Boolean object on which the `not` method is invoked, e.g. `ready.not`.

<sub>*Postfix `++` and `--` are not legal in **dish**</sub>

## Syntax
**dish** has 4 types of syntax blocks, Scalar Blocks `(...)`, Array blocks `[...]`, Hash blocks `{...}`, and Logic blocks. Expressions and statements in each block are terminated either by a statement ending newline<sup>*</sup>, by a comma, or by a semicolon (unless the newline, comma, or semicolon is encapsulated in a string literal).

Scalar blocks `(...)` return the last expression or statement in the block. A full **dish** program is inside an implicit Scalar block when the first and last characters of the program are not `(` and `)` respectively. Array blocks `[...]` and Hash blocks `{...}` return the entire array or hash, where Hash blocks contain all locally-scoped variables when returned--allowing additional logic to exist in Hash block intitialization without changing its structure.

Logic blocks are represented by a colonized list of arguments followed by a Scalar block, as in `:x(...)` or `:x:y(...)`, returning a data type corresponding to the last expression. A Logic block may have no arguments, but still must be preceded by a single colon character as in `:(...)`. All arguments passed to a Logic block are locally scoped.

If a Logic block conatins only a single argument and has been invoked on a Hash or an Array, then the entire Hash or Array is passed as that argument; when the Logic block contains multiple arguments and has been invoked on a Hash or an Array, then each Hash or Array item are passed as the arguments.

Comments in **dish** start with a double pound `##` and end with a newline. There are no multi- or in-line comments in **dish**.

<sub>*Statement ending newlines are those not preceded by an opening block character or by a binary operator.</sub>

## Variables and Member Access
Variables in **dish** must start with a dollar sign or a letter, followed by any number of numbers, letters, or underscores.

Member access in **dish** is indicated by the special `.` operator which precedes a member expression. That expression is evaluated, and the member at that evaluated expression is returned. Each "member" of a String is indexed numerically and represents the character at that index; Each "member" of a Number is indexed numerically and represents each bit, with the least significant bit at index 0. Note that to retrieve the member index of a static Number, parenthesis are often necessary to disambiguate Numeric member access from a floating point. Compare the following:

`dish -e '12.1'` outputs the floating point number: `12.1`

`dish -e '12.(1)' outputs the 0th-indexed 2nd most-significant bit of the integer 12: `1`

`dish -e 'a = 12; a.1'` outputs the 0th-indexed 2nd most-significant bit (disambiguated by the variable name `a`): `1`

Because **dish** variables can not begin with numbers, numeric index members such as `[1, 2, 3].1` or `'foobar'.3` are unambiguous; Hash key members, however, because they can be ambiguous, should be quoted to avoid already-defined variable names. Compare the following:

`dish -e 'a = {foo: 1, bar: 2}; foo = "bar"; a."foo"'` outputs `1`

`dish -e 'a = {foo: 1, bar: 2}; foo = "bar"; a.foo'` outputs `2`

The member expression can also be a Logic block, as in `[1, 2, 3].:a:b:c(a + b + c)` or a variable that points to a Logic block as in `power = :a:b(a ^ b); [2, 3].power`. As seen in these last two examples, the values of a List data type are passed as arguments to the Logic block. This is similar for Scalar data types as in:

`dish -e 'squared = :n(n ^ 2); 3.squared'`

`dish -e 'title = :s(s.words.map:w(w.0 @= w.0.uc).join(" ")); "my title".title'`

In cases where a Logic block contains exactly two arguments, you can use the following, alternative binary syntax for passing arguments: `power = :a:b(a ^ b); 2.power(3)` where the first argument is the object on which the Logic block is invoked and where the second argument is passed from a Scalar block following the member expression. Similarly, in cases where a Logic block contains more than two arguments, you can use the following, alternative n-ary syntax for passing arguments: `quad = :x:a:b:c(a * x ^ 2 + b * x + c); 2.quad[2, 3, 4]`, where additional arguments are passed from an Array block following the member expression. Note, however, that `2.quad(2, 3, 4)` would only pass `4` as a second argument since `(2, 3, 4)` is a Scalar block that only returns the last expression.

## String Interpolation
**dish** supports string interpolation by injecting a Scalar block prefixed with a `$` character inside a double-quoted string (`"$(...)"`), for example: `dish -e '(0..9).map:i("i: $(i)").join("\n")'`. Any **dish** expression can be included in a string interpolated Scalar block, but you will need to escape any double quote characters used in your expression.

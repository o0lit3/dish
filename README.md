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

## Output
By default, **dish** outputs the last evaluated statement to STDOUT. If the last evaluated statement is an Array or a Hash, the output is formatted as valid JSON. If the last evaluated statement is a Scalar, the scalar's raw output is printed to STDOUT.

This behavior allows you to pipe the output of one dish executable into another and use it as JSON input via `stdin` downstream.

## Data Types and Operators
For a detailed decription of operators, precedence, and implicit operator context, [read the operator documentation at tests/README.md](tests/README.md).

Like JSON itself, **dish** has 6 data types (Hash, Array, String, Number, Boolean, and Null), all of which are objects with their own object methods. Traditional, symbolic operators (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods.

As such, all binary operators apply the right-hand operand as a parameter to a method on the left-hand object. In the expression `5 + 2`, `2` is a parameter of the `add` method on the Number object `5`, e.g. `5.add(2)`.

Similarly, traditional unary operators (which are always prefix operators when represented in shorthand notation<sup>*</sup>) correspond to methods with no parameters that are invoked on the single operand as an object. In the expression `!ready`, `ready` is the Boolean object on which the `not` method is invoked, e.g. `ready.not`.

<sub>*Postfix `++` and `--` are not legal in **dish**</sub>

## Syntax
**dish** has 4 types of syntax blocks, Scalar Blocks `(...)`, Array blocks `[...]`, Hash blocks `{...}`, and logic blocks. Expressions and statements in each block are terminated either by a statement ending newline<sup>*</sup>, by a comma, or by a semicolon (unless the newline, comma, or semicolon is encapsulated in a string literal).

Scalar blocks `(...)` return the last expression or statement in the block. A full **dish** program is inside an implicit scalar block when the first and last characters of the program are not `(` and `)` respectively. Array blocks `[...]` and Hash blocks `{...}` return the entire array or hash, where hash blocks contain all locally-scoped variables when returned.

Logic blocks are represented by a colonized list of arguments followed by a Scalar block, as in `:x(...)` or `:x:y(...)`, returning a data type corresponding to the last expression. A logic block may have no arguments, but still must be preceded by a single colon character as in `:(...)`. All arguments passed to a logic block are locally scoped.

Comments in **dish** start with a double slash `//` and end with a newline. There are no multi- or in-line comments in **dish**.

<sub>*Statement ending newlines are those not preceded by an opening block character or by a binary operator.</sub>

## Variables and Member Access
Variables in **dish** must start with either a letter or a `$` character, and may only contain letters, numbers, underscores (`_`), and dollar signs (`$`).

Member access in **dish** is indicated by the special `.` operator which precedes a member expression. That expression is evaluated, and the member at that evaluated expression is returned. Because **dish** variables can not begin with numbers, array index members such as `a = [1, 2, 3]; a.0` are unambiguous; Hash key members, however, because they can be ambiguous, should be quoted. Compare the following:

`dish -e 'a = {foo: 1, bar: 2}; foo = "bar"; a."foo"'` outputs `1`

`dish -e 'a = {foo: 1, bar: 2}; foo = "bar"; a.foo'` outputs `2`

The member expression can also be a Logic block, as in `[1, 2, 3].:a:b:c(a + b + c)` or a variable that points to a Logic block as in `power = :a:b(a ** b); [2, 3].power`. As seen in these last two examples, the values of a List data type are passed as arguments to the Logic block. This is similar for Scalar data types as in `"foobar".:s(s.upper)` or `squared = :a(a ** 2); 3.squared`.

In cases where a Logic block contains exactly two arguments, you can use the the following, alternative binary syntax for passing arguments: `power = :a:b(a ** b); 2.power(3)` where the first argument is the object on which the Logic block is invoked and where the second argument is passed from a Scalar block following the member expression.

## String Interpolation
**dish** supports string interpolation by injecting a Scalar block prefixed with a `$` character inside a double-quoted string (`"$(...)"`), for example: `dish -e '(0..9).map:i("i: $(i)").join("\n")'`. Any **dish** expression can be included in a string interpolated Scalar block, but you will need to escape any double quote characters used in your expression.


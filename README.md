# dish
Dish is a dynamically-typed, keywordless programming language designed with the following goals:

* Allow quick and easy processing of json data with traditional programming operators (as opposed to the query-based operators of [jq](https://github.com/stedolan/jq#readme))
* Provide flexibility in language syntax so that developers can tackle tasks based on individual conceptualization (while still allowing peers to easily follow logical flow--that is, while avoiding the "write-only" criticsms of [perl](https://github.com/Perl/perl5#readme))
* Isolate behavioral functionality by treating all data types as objects (and remove control flow keywords in favor of object methods, which goes a step further than [ruby](https://github.com/ruby/ruby#readme))

## Installation
Dish is interpreted by [Go](https://github.com/golang/go#readme). With Go installed, build the **dish** interpreter via `go build -o /usr/local/bin/dish` from the project root. You can then run **dish** files via `dish /path/to/file.dish` or with the `-e` command flag, as in `dish -e '["Hello World!"]'"`

## Output
By default **dish** outputs the last evaluated statement to STDOUT. If the last evlauated statement is an [Array](#array) or a [Hash](#hash), the output is formatted as valid JSON. If the last evaluated statement is a Scalar, the scalar's raw output printed to STDOUT.

## Data Types
Dish has 5 data types, all of which are objects with their own object methods. Traditional, symbolic operands (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods.

As such, all binary operators apply the right-hand operand as a parameter to a method on the left-hand object. In the expression `5 + 2`, `2` is a parameter of the `plus` method on the Number object `5`, e.g. `5.plus(2)`.

Similarly, traditional unary operators (which must be used as prefix operators when represented in shorthand notation; postfix `++` and `--` are not legal in **dish**) correspond to methods with no parameters that are invoked on the single operand as an object. In the expression `!ready`, `ready` is the Boolean object on which the `not` method is invoked, e.g. `ready.not()` 

### Boolean

### Number

### String

### Array

### Hash

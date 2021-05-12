# dish
Dish is a dynamically-typed programming language designed with the following goals:

* Allow quick and easy processing of json data with traditional programming operators (as opposed to the query-based operators of [jq](https://github.com/stedolan/jq#readme))
* Provide flexibility in language syntax so that developers can tackle tasks based on individual conceptualization (while still allowing peers to easily follow logical flow--that is, while avoiding the "write-only" criticsms of [perl](https://github.com/Perl/perl5#readme))
* Isolate behavioral functionality by treating all data types as objects (and remove control flow keywords in favor of object methods, which goes a step further than [ruby](https://github.com/ruby/ruby#readme))

## Data Types
Dish has 5 data types, all of which are objects with their own object methods. Traditional, symbolic operands (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods. 

### Boolean

### Number

### String

### Array

### Hash

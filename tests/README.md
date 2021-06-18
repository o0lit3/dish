# Operators
Symbolic operators in *dish* (like `+`, `-`, `*`, `/`, etc.) are shorthand representations of longer form object methods. The tables below outline which shorthand, symbolic operands correspond to which method name. Some binary operators change behavior depending on the data types of their operands; in these cases, the left-hand operand forces the right-hand operand into the same data type context with the following conversion rules:

|         | Boolean             | Number  | String            | Array    | Hash    |
| ------- | ------------------- | ------- | ----------------- | -------- | ------- |
| Boolean | x                   | 1, 0    | "true", "false"   | [x]      | {~x: x} |
| Number  | x != 0              | x       | ~x                | [x]      | {~x: x} |
| String  | x != "0" && x != "" | +x || 0 | x                 | [x]      | {~x: x} |
| Array   | #x != 0             | #x      | x.join('')        | x        | x.hash  |
| Hash    | #x != 0             | #x      | x.values.join('') | x.values | x       |

With the exception of special assignment operators, all symbolic operators in *dish* consist of either a single symbolic character or a doubled, symbolic character. Restated: no symbolic operator in dish (except the specialized [assigment operators](#assignment-operators)) contain a mixture of two different symbolic characters, and all symbolic operators are either 1 or 2 characters long.

## Arithmetic Operators
| Operator | Operands                          | Method Name   | Example                            |
| -------- | --------------------------------- | ------------- | ---------------------------------- |
| +        | Scalar + Implicit Numeric Context | a.add(b)      | 10 + 20 = 30                       |
|          | Array + Implicit Array Context    | a.concat(b)   | [1, 2] + [2, 3] = [1, 2, 2, 3]     |
|          | Hash + Implicit Hash Context      | a.concat(b)   | {x: 1} + {y: 2} = {"x": 1, "y": 2} |
|          |                                   |               |                                    |
| -        | Scalar + Implicit Numeric Context | a.subtract(b) | 10 - 20 = -10                      |
|          | Array + Implicit Array Context    | a.remove(b)   | [1, 2] - [2, 3] = [1]              |
|          | Hash + Implicit Hash Context      | a.remove(b)   | {x: 1, y: 2} - {"x"} = {"y": 2}    |

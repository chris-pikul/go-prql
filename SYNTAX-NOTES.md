# Syntax and Parsing Notes

This document holds some notes, tips, and documentation on the PRQL syntax and
some RFC comments on ways to improve the parsing.

## Table of Contents

- [Syntax - General lexiconal grammar and explanation](#syntax)
- [Header Directive](#header-directive)
  - [Dialect](#dialect)
  - [Version](#version)
- [Functions](#functions)
  - [Named-Parameters](#named-parameters)
  - [Positional-Parameters](#positional-parameters)
  - [Function Body](#function-body)
  - [Usage - Calling/Invocing](#function-usage)
    - [Explicit Invocation](#explicit-invocation)
    - [Implied Invocation](#implicit-invocation)
  - [Function as a Constant](#function-as-a-constant)

## Syntax

The basic syntax of PRQL follows _statements_ which are delimeted by either an
EOL (new line), or the _pipe operator_ `|`. Within the official docs, referring
to the [PRQL Language Book](https://prql-lang.org/book/): _transforms_
are _statements_, and the query is considered a _pipeline_ of _transforms_ as a
whole.

Following this, the general syntax of PRQL can be considered
```
transform ::== aggregate | derive | filter | from | group | join | select | sort | take | window
keyword ::== prql | func | table | {transform}
reference ::== {literal} | {assignment}

comparison_operator ::== > | >= | < | <= | == | !=
logical_operator ::== and | or

boolean_expression ::== {identifier} {comparison_operator} {identifier | literal}
assignment ::== {identifier} = {expression}
range ::== ( {identifier} | in {integer}{? .. {integer}} )

named_var ::== {identifier}:{literal}

statement ::== {keyword} {expression}
pipeline ::== {statement} {? | statement... }
```

Here we have a set of known _transform keywords_ which declare at the start of
the _statement_ the intent of the statement. These are required at the beginning
of the statement line and declare the context for the remaining parsing.

The following is a table of known _transform keywords_ and their individual
syntaxs using the formating declared above:

| Transform Keyword | Purpose | Syntax |
|-------------------|---------|--------|
| aggregate | Convert many rows into a singular row | `aggregate [ {expression|assignment},... ]` |
| derive | Compute new columns | `derive [ {assignment}... {,} ]` |
| filter | Pick rows by value | `filter {boolean_expression {? {logical_operator} {boolean_expression}}... | {range} }` |
| from | Specifies the data source | `from {reference}` |
| group | Partitions rows into groups with pipelines applied | `group {literal | [ {literal},... ]} {pipeline}` |
| join | Adds columns from another table, matching on condition | `join side:{inner|left|right|full} {literal} {[ {boolean_expression},... ]}` |
| select | Picks and computes columns | `select {{assignment} | [ {assignment},... ]}` |
| sort | Orders rows based on columns | `sort {{+|-}{literal} | [ {{+|-}{literal}},... ]}` |
| take | Pick rows based on position | `take {{number} | {range}}` |
| window | Applies a pipeline to segments of rows | _see book_ |

## Header Directive

The specification dictates that a _query_ or _document_ may begin with a special
directive `prql` declaring the SQL dialect and PRQL version. This directive
**MUST** be present at the beginning of the document. A _document_ is defined as
the string (contents) of PRQL source comprising any functions, tables, and
pipelines.

```
prql {? dialect:{ENUM}} {? version:{INTEGER}}
```

Here we see that the keyword starting the header directive starts with `prql`
and can feature two _named-variables_. Each parameter is considered optional,
and in the case of not being provided will assume reasonable defaults as noted
bellow. This syntax follows [Function](#functions) usage.

Because each parameter is optional, no arguments are required to be provided to
the `prql` statement, and a blank directive is possible (although useless).

Examples:
```PRQL
prql dialect:pgsql
prql version:1
prql dialect:mysql version:1
prql
```

### Dialect

The `dialect` parameter takes an enumeration as it's argument. This parameter
defaults to `generic` when not supplied.

The supported string-enum values are:

- `ansi`
- `bigquery`
- `clickhouse`
- `generic` (default)
- `hive`
- `mssql`
- `mysql`
- `postgres`
- `sqlite`
- `snowflake`

> **NOTE:** The specification does not specify how to handle invalid/unknown
> dialect options. As such, this implementation assumes any unknowns are
> recoverable errors (warnings) and will be assumed to default as `generic`.

### Version

The `version` parameter takes an number argument and declares the PRQL syntax
version that the remaining document will use. It defaults to `1`, being the
current version.

> **NOTE:** The specification does not specify the format of acceptable
> version arguments. As such, this implementation assumes it to be an
> unsigned-integer. This is based on the compiler source code assigning an `i64`
> to this value.
>
> @see [prql/prql-compiler/src/ast/query.rs#L8](https://github.com/prql/prql/blob/2cd32b9d1fd6e4e8f58d8351ddb8cb7a2012e41e/prql-compiler/src/ast/query.rs#L8)

## Functions

PRQL is extensible through the use of functions. These functions are included in
the _document_ and can be referenced within the query pipelines. Functions
should be declared before they are used (spec doesn't specify).

```
positional_param ::== {identifier}
named_param ::== {identifier}:{literal}

func {identifier} {named_param...} {positional_param...} -> {expression}
```

The first identifier is the given name (or symbol) representing the function for
use in later pipelines. When using functions, this identifier is used to call
this function.

Following the identifier is any _named-parameters_ if any.
Following the _named-parameters_ is any _positional-parameters_ if any.

After the function declaration is the arrow operator `->` signifying the start
of the function body. Being that function bodies are just expressions they adhere
to the same rules as other expressions and precedence.

Examples:
```PRQL
func pi -> 3.14159
func deg_to_rad deg -> deg * (pi / 180)
func interpolate low:0 high val -> (val - low) / (high - low)
```

### Named-Parameters

```
named_parameter ::== {identifier}:{literal}
```

Following the function name identifier, and **before** any _positional-parameters_,
are any _named-parameters_ declarations.

There can be 1 or more named-parameters which are space delimited. All named
parameters must be together, meaning that named-parameters and
positional-parameters cannot be mixed (interweaved).

Each named-parameter declaration follows the format of the parameter identifier
being first, followed by a colon operator `:`, and finalized with a default
value literal.

> **Every named-parameter MUST be declared with a default value**

When calling functions, named-parameters are considered optional parameters due
to their default values being declared in the function declaration.

### Positional-Parameters

```
positional_param ::== {identifier}
```

At the end of the function definition, **before** the start of the function body
(operator), can be any _positional-parameters_.

They are declared only by an identifier which can be referenced within the
function body scope. There can be 1 or more positional-parameters which are
space delimited. All positional-parameters must be together, meaning that
named-parameters and positional-parameters cannot be mixed (interweaved).

These parameters are considered positional in that when calling a function the
identifiers are not specified and instead the position of any given arguments
are mapped to the matching index on the function declaration.

> **Every positional-parameter IS REQUIRED when calling the function**

### Function Body

After the function declaration (identifier, and parameters) the function
definition, or function body, can be declared. This is started by the arrow
operator `->`.

The actual body of a function is an expression that **MUST** resolve to some
value. Functions almost act as macros in other languages like C++ in that the
body effectively has any parameters replaced with the arguments given during
calling and then simplified if necessary and piped directly into the query.

> **NOTE:** The specification does not directly specify this is how functions
> work, this statement is inferred from usage with the PRQL Playground and how
> output SQL seems to be generated.

Given that the function body resolves to an expression, you cannot perform any
transformations or sub-queries within it. It is only used for deriving values
through calculation expressions.

Functions **cannot** be recursive. That is, a function cannot call itself, being
that the function identifier (symbol) is not declared in scope until after the
function declaration.

Functions **can** call other functions that are declared before them, allowing
chaining of functions and extended compositions.

### Function Usage

To use, or call a function, only the identifier (symbol) is required. With this,
there are two ways to use a function. Either called directly, or used with a
_pipe_ operator `|`. Called directly is considered _explicit invocation_, while
used with a pipe is considered _implicit invocation_.

Regardless of the use, any arguments provided to parameters, either named or
positional can be provided as _expressions_ which will be substituted directly
when compiled. Meaning, that an argument can be either a literal, a reference
such as column name, or an expression using any combination of the prior that
results in a singular value.

#### Explicit Invocation

When called directly, or _explicit invocation_, the syntax follows the following:

```
explicit_invoc ::== {identifier} {? {named_param}... } {? {positional_param}... }
```

An explicit invocation **MUST** be wrapped in paranthesis for precedence, this
allows the parser to properly tokenize the parameters.

Named-parameters must be specified first if they are desired (differ from the
defaults), followed by any positional arguments as values themselves.

Example:
```
func celcius_to_fahrenheit deg_c -> (deg_c * (9/5)) + 32

from weather
derive deg_f = (celcius_to_fahrenheit deg_c)
```

#### Implicit Invocation

When called as part of a pipeline, using the pipe operator `|`, the syntax
remains the same as an _explicit invocation_ with the major difference that the
previous value being "piped" into this function call will take the place as the
**last positional parameter**. It is this functionality of substituting the last
positional parameter with the current pipeline value (state) that gives this
call style the name "implicit invocation".

> **NOTE:** The specification does not declare any requirement for the function
> to have any positional-parameters. In fact, since functions act almost as
> macros, substituting variables within expressions, a _constant function_ will
> just replace the output of the pipe with it's value with no warnings given.

#### Function as a Constant

Since all parameters are optional in the declaration of the function, a function
declared with no parameters can act as a constant variable, such as:

```PRQL
func pi -> 3.14159

from pizzas
derive volume = ((radius * radius) * pi) * thickness
```

The function above when called will effectively act as a constant within your
queries replacing the function call with the literal directly. Resulting in the
following SQL output:

```SQL
SELECT
  pizzas.*,
  radius * radius * 3.14159 * thickness AS volume
FROM
  pizzas
```

As noted in the [_implicit invocation_](#implicit-invocation) sub-chapter, you
can pipe into a contant function such as above. Since the function takes no
arguments the previous state of the pipeline is discarded and replaced with the
constant value returned. This functionality may be a "dragon".

# Syntax and Parsing Notes

This document holds some notes, tips, and documentation on the PRQL syntax and
some RFC comments on ways to improve the parsing.

## Syntax

The basic syntax of PRQL follows _statements_ which are delimeted by either an
EOL (new line), or the _pipe operator_ `|`. Within the official docs, referring
to the [PRQL Language Book](https://prql-lang.org/book/): _transforms_
are _statements_, and the query is considered a _pipeline_ of _transforms_ as a
whole.

Following this, the general syntax of PRQL can be considered
```
transform ::== aggregate | derive | filter | from | group | join | select | sort | take | window
keyword ::== func | table | {transform}
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
| window | Applies a pipeline to segments of rows | _see below_ |


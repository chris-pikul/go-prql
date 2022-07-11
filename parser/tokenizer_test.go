package parser

import (
	"fmt"
	"testing"
)

func TestTokenizer(t *testing.T) {
	const input = `from emp = employees
	filter country_code == "USA"   # Each line transforms the previous result.
	derive [                       # This adds columns / variables.
		gross_salary = s'salary + payroll_tax',
		gross_cost = gross_salary + benefits_cost  # Variables can use other variables.
	]
	filter gross_cost > 0
	group [ title, country_code] (  # For each group use a nested pipeline
		aggregate [                  # Aggregate each group to a single row
			average salary,
			average gross_salary,
			sum salary,
			sum gross_salary,
			average gross_cost,
			sum_gross_cost = sum gross_cost,
			ct = count,
		]
	)
	sort sum_gross_cost
	filter ct > 200 | take 20
	join countries side:left [country_code]
	derive [
		always_true = true,
		db_version = s''''version()'''',    # An S-string, which transpiles directly into SQL
	]`

	tokens := tokenize(input)
	fmt.Printf("Tokenized, generated %d tokens\n", len(tokens))
	for ind, tkn := range tokens {
		fmt.Printf("[%d] %d:%d %s `%s`\n", ind, tkn.Line, tkn.Character, tkn.Type.String(), tkn.Value)
	}

	t.Fail()
}

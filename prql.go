package prql

// Compile takes an incoming PRQL query (string) and returns the SQL standard
// equivelent (string), or an error if one occured. In the event of an error,
// the string returned will be empty. The error type is a custom type wrapping
// the go standard error as "prql.Error".
func Compile(source string) (string, *Error) {
	return "", nil
}

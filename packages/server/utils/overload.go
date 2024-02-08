package utils

// If the condition is true, return the trueValue, otherwise return the falseValue.
func TernaryOperator(condition bool, trueValue interface{}, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}

// If the condition is false, panic with the message
func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

// Return true if the string e is in the array s, otherwise return false.
func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// It takes a variable number of maps and returns a single map that contains all the keys and values
// from the input maps
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Pair is a struct with two fields, First and Second, of types T and U, respectively.
// @property {T} First - The first element of the pair.
// @property {U} Second - U
type Pair[T, U any] struct {
	First  T
	Second U
}

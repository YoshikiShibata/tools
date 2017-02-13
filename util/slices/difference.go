// Copyright © 2017 Yoshiki Shibata. All rights reserved.

package slices

// DifferenceString returns the difference between two slices: x - y.
func DifferenceString(x []string, y []string) []string {
	if len(x) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(x))

	if len(y) == 0 {
		result = append(result, x...)
		return result
	}

	for _, v := range x {
		if !ContainsString(y, v) {
			result = append(result, v)
		}
	}

	return result
}

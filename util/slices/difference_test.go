// Copyright © 2017 Yoshiki Shibata. All rights reserved.

package slices

import (
	"reflect"
	"testing"
)

func TestDifferenceString(t *testing.T) {
	for _, test := range []struct {
		x        []string
		y        []string
		expected []string
	}{
		{nil, nil, []string{}},
		{[]string{"hello"}, nil, []string{"hello"}},
		{[]string{"hello"}, []string{"hello"}, []string{}},
		{[]string{"hello", "world"}, []string{"hello"}, []string{"world"}},
		{[]string{"hello", "world"}, []string{"world"}, []string{"hello"}},
		{[]string{"hello", "world"}, []string{"hello", "world"}, []string{}},
	} {
		result := DifferenceString(test.x, test.y)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("difference between %v and %v is %v, but want %v\n",
				test.x, test.y, result, test.expected)
		}
	}
}

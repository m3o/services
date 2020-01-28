package unique

import "testing"

func TestStrings(t *testing.T) {
	tt := []struct {
		name   string
		input  []string
		output int
	}{
		{name: "All unique", input: []string{"foo", "bar", "dog"}, output: 3},
		{name: "Some unique", input: []string{"foo", "bar", "bar"}, output: 2},
		{name: "One unique", input: []string{"foo", "foo", "foo"}, output: 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if res := Strings(tc.input); len(res) != tc.output {
				t.Errorf("Unexpected error, got %v, expected %v", len(res), tc.output)
			}
		})
	}
}

package patch

import (
	"testing"
)

func TestMergeTags(t *testing.T) {
	tests := []struct {
		name           string
		oldTag, newTag string
		want           string
	}{
		{
			name:   "NotOverrideTag",
			oldTag: `json:"code;omitempty"`,
			newTag: `test:"test"`,
			want:   `json:"code;omitempty" test:"test"`,
		},
		{
			name:   "OverrideSingleTag",
			oldTag: `json:"code;omitempty"`,
			newTag: `json:"-" test:"test"`,
			want:   `json:"-" test:"test"`,
		},
		{
			name:   "OverrideMultiTag",
			oldTag: `json:"code;omitempty" test1:"test"`,
			newTag: `test1:"test1" test2:"test2"`,
			want:   `json:"code;omitempty" test1:"test1" test2:"test2"`,
		},
		{
			name:   "EmptyValueTag",
			oldTag: `json:"code;omitempty" test1 test2:"test2"`,
			newTag: `test2`,
			want:   `json:"code;omitempty" test2:"test2"`,
		},
		{
			name:   "WithColonTag",
			oldTag: `json:"code;omitempty"`,
			newTag: `test1:"a:b" test2:"test2"`,
			want:   `json:"code;omitempty" test1:"a:b" test2:"test2"`,
		},
		{
			name:   "WithSpacesTag",
			oldTag: `json:"code;omitempty"`,
			newTag: `test1:"a b" test2:"test2"`,
			want:   `json:"code;omitempty" test1:"a b" test2:"test2"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := mergeTags("`"+test.oldTag+"`", "`"+test.newTag+"`"); got != "`"+test.want+"`" {
				t.Fatalf(" got: %s\nwant: %s\n", got, test.want)
			}
		})
	}
}

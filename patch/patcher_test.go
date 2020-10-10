package patch

import (
	"testing"
)

func TestMergeTags(t *testing.T) {
	tests := []struct {
		name string
		in   [2]string
		want string
	}{
		{
			name: "NotOverrideTag",
			in: [2]string{
				"`json:\"code;omitempty\"`",
				"`test:\"test\"`",
			},
			want: "`json:\"code;omitempty\" test:\"test\"`",
		},
		{
			name: "OverrideSingleTag",
			in: [2]string{
				"`json:\"code;omitempty\"`",
				"`json:\"-\" test:\"test\"`",
			},
			want: "`json:\"-\" test:\"test\"`",
		},
		{
			name: "OverrideMultiTag",
			in: [2]string{
				"`json:\"code;omitempty\" test1:\"test\"`",
				"`test1:\"test1\" test2:\"test2\"`",
			},
			want: "`json:\"code;omitempty\" test1:\"test1\" test2:\"test2\"`",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := mergeTags(test.in[0], test.in[1]); got != test.want {
				t.Fatalf(" got: %s\nwant: %s\n", got, test.want)
			}
		})
	}
}

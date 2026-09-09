package greeting

import "testing"

func TestSquare(t *testing.T) {
	cases := map[string]struct {
		in   int
		want int
	}{
		"zero":     {0, 0},
		"positive": {7, 49},
		"negative": {-4, 16},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Square(tc.in); got != tc.want {
				t.Fatalf("Square(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

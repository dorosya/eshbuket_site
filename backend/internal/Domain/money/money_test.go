package money

import "testing"

func TestParseToCents_Success(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{in: "100", want: 10000},
		{in: "12.5", want: 1250},
		{in: "12.50", want: 1250},
		{in: "12,50", want: 1250},
		{in: " 1 ", want: 100},
	}

	for _, tc := range cases {
		got, err := ParseToCents(tc.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("unexpected cents for %q: got %d want %d", tc.in, got, tc.want)
		}
	}
}

func TestParseToCents_Invalid(t *testing.T) {
	cases := []string{"", "abc", "12.345", "-1", "0", ".", "1.", "1..2"}

	for _, in := range cases {
		if _, err := ParseToCents(in); err == nil {
			t.Fatalf("expected error for %q", in)
		}
	}
}

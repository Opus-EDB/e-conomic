package economic

import "testing"

func TestNormalizeCorporateId(t *testing.T) {
	cases := map[string]string{
		"":               "",
		"0":              "",
		"000":            "",
		"25781309":       "25781309",
		"0025781309":     "25781309", // order 26102: typed with leading zeros
		"DK 0025781309":  "25781309",
		"0 0 25781309":   "25781309",
		".10521815":      "10521815",
		"DK 10 52 18 15": "10521815",
		"10.521.815":     "10521815",
		"10500100":       "10500100", // inner and trailing zeros stay
	}
	for in, want := range cases {
		if got := NormalizeCorporateId(in); got != want {
			t.Errorf("NormalizeCorporateId(%q) = %q, want %q", in, got, want)
		}
	}
}

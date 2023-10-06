package utils

import (
	"testing"
)

func tokTry(t *testing.T, str string, tokens []string, ref []string) {
	t.Logf("starting split: %s", str)
	ret := SplitTokens(str, tokens)
	t.Logf("splitted: %v", ret)

	for x := range ref {
		if ret[x] != ref[x] {
			t.Fatalf("token #%d mismatch: is \"%s\" should be \"%s\"", x, ref[x], ret[x])
		}
	}
}

func TestTokenize(t *testing.T) {
	tokens := []string{"${", "}", "$(", ")"}

	tokTry(t, "foo ${one::two} abc", tokens, []string{"foo ", "${", "one::two", "}", " abc"})
	tokTry(t, "${a::b}", tokens, []string{"${", "a::b", "}"})
	tokTry(t, "$(cap a, b, c)", tokens, []string{"$(", "cap a, b, c", ")"})
}

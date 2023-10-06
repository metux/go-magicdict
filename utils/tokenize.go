package utils

import (
	"strings"
)

func SplitTokens(s string, tokens []string) []string {
	res := make([]string, 0)

	pos := 0
	for len(s) > pos {
		tail := s[pos:]
		move := true

		for _, tok := range tokens {
			if strings.HasPrefix(tail, tok) {
				if pos > 0 {
					res = append(res, s[:pos])
				}
				res = append(res, tok)
				s = s[(pos + len(tok)):]
				pos = 0
				move = false
			}
		}
		if move {
			pos = pos + 1
		}
	}

	if len(s) > 0 {
		res = append(res, s)
	}
	return res
}

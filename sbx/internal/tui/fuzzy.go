package tui

import (
	"strings"
	"unicode"
)

// fzf-style scoring: every query character must appear in order, and the
// alignment with the most word-boundary and consecutive hits wins.
const (
	scoreMatch       = 16
	scoreGapStart    = -3
	scoreGapExtend   = -1
	bonusBoundary    = 8
	bonusCamel       = 7
	bonusConsecutive = 4
	firstCharFactor  = 2
)

// fuzzyScore matches each space-separated term of q against s independently
// (like fzf's extended mode) and returns the summed score. ok is false when
// any term fails to match. Matching is case-insensitive.
func fuzzyScore(s, q string) (score int, ok bool) {
	for _, term := range strings.Fields(q) {
		v, matched := termScore(s, term)
		if !matched {
			return 0, false
		}
		score += v
	}
	return score, true
}

func charBonus(prev, cur rune) int {
	switch {
	case prev == 0 || !unicode.IsLetter(prev) && !unicode.IsDigit(prev):
		return bonusBoundary
	case unicode.IsLower(prev) && unicode.IsUpper(cur), !unicode.IsDigit(prev) && unicode.IsDigit(cur):
		return bonusCamel
	}
	return 0
}

func termScore(s, term string) (int, bool) {
	text := []rune(s)
	pattern := []rune(strings.ToLower(term))
	n, m := len(text), len(pattern)
	if m == 0 {
		return 0, true
	}
	if m > n {
		return 0, false
	}
	lower := make([]rune, n)
	bonus := make([]int, n)
	for j := range text {
		lower[j] = unicode.ToLower(text[j])
		var prev rune
		if j > 0 {
			prev = text[j-1]
		}
		bonus[j] = charBonus(prev, text[j])
	}
	const none = -1 << 30
	// prev[j]: best score with the previous pattern char matched at text[j].
	prev := make([]int, n)
	for j := range n {
		prev[j] = none
		if lower[j] == pattern[0] {
			prev[j] = scoreMatch + bonus[j]*firstCharFactor
		}
	}
	for i := 1; i < m; i++ {
		cur := make([]int, n)
		for j := range n {
			cur[j] = none
			if lower[j] != pattern[i] {
				continue
			}
			for k := j - 1; k >= 0; k-- {
				if prev[k] == none {
					continue
				}
				v := prev[k] + scoreMatch + bonus[j]
				if gap := j - k - 1; gap == 0 {
					v += max(bonusConsecutive, bonus[k])
				} else {
					v += scoreGapStart + scoreGapExtend*(gap-1)
				}
				cur[j] = max(cur[j], v)
			}
		}
		prev = cur
	}
	best := none
	for _, v := range prev {
		best = max(best, v)
	}
	return best, best != none
}

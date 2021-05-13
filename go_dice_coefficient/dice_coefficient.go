package main

import "sort"

type bigram struct {
	digit1 rune
	digit2 rune
}

func (b bigram) equals(other bigram) bool {
	return b.digit1 == other.digit1 && b.digit2 == other.digit2
}

func (b bigram) lessThan(other bigram) bool {
	return b.digit1 < other.digit1 ||
		(b.digit1 == other.digit1 && b.digit2 < other.digit2)
}

func CalculateDiceCoefficient(text1 string, text2 string) float64 {
	if text1 == `` || text2 == `` {
		return 0
	}
	if text1 == text2 {
		return 1
	}
	if len(text1) == 0 || len(text2) == 0 {
		return 0
	}
	start := 0
	runes1 := []rune(text1)
	bigrams1 := make([]bigram, len(runes1)-1)
	for end := 1; end < len(runes1); end++ {
		bigrams1[start] = bigram{runes1[start], runes1[end]}
		start++
	}
	sort.Slice(bigrams1, func(i int, j int) bool {
		return bigrams1[i].lessThan(bigrams1[j])
	})

	start = 0
	runes2 := []rune(text2)
	bigrams2 := make([]bigram, len(runes2)-1)
	for end := 1; end < len(runes2); end++ {
		bigrams2[start] = bigram{runes2[start], runes2[end]}
		start++
	}
	sort.Slice(bigrams2, func(i int, j int) bool {
		return bigrams2[i].lessThan(bigrams2[j])
	})

	i1 := 0
	i2 := 0
	matches := 0

	for i1 < len(bigrams1) && i2 < len(bigrams2) {
		if bigrams1[i1].equals(bigrams2[i2]) {
			matches++
			i1++
			i2++
			continue
		}
		if bigrams1[i1].lessThan(bigrams2[i2]) {
			i1++
			continue
		}
		i2++
	}

	return float64(matches*2) / float64(len(bigrams1)+len(bigrams2))
}

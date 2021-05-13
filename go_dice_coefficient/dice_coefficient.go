package main

func CalculateDiceCoefficient(text1 string, text2 string) float64 {
	if text1 == `` || text2 == `` {
		return 0
	}
	if text1 == text2 {
		return 1
	}
	return 0
}

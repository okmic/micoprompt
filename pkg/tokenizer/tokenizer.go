package tokenizer

import "unicode/utf8"

func Estimate(text string) int {
	runes := utf8.RuneCountInString(text)
	if runes == 0 {
		return 0
	}
	nonASCII := 0
	for _, r := range text {
		if r > 127 {
			nonASCII++
		}
	}
	ratio := 4.0
	if float64(nonASCII)/float64(runes) > 0.2 {
		ratio = 3.0
	}
	return int(float64(runes) / ratio)
}

package tokenizer

import "unicode/utf8"

// Estimate возвращает примерное количество токенов.
// Эвристика: для английского ~4 символа на токен, для кода/русского ~3.
func Estimate(text string) int {
	runes := utf8.RuneCountInString(text)
	if runes == 0 {
		return 0
	}
	// считаем долю не-ASCII
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

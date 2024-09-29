package handlers

import (
	"strings"
	"unicode"
)

func replaceWords(input string, wordsToReplace []string) string {
	words := strings.Fields(input)
	for i, word := range words {
		cleanWord := strings.TrimFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r)
		})

		for _, target := range wordsToReplace {
			if strings.EqualFold(cleanWord, target) {
				words[i] = strings.Repeat("*", 4) + word[len(cleanWord):]
				break
			}
		}
	}
	return strings.Join(words, " ")
}

func validateChirp(input string) string {
	wordsToReplace := []string{"kerfuffle", "sharbert", "fornax"}
	return replaceWords(input, wordsToReplace)
}

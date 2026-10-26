package helper

import (
	"strings"
)

var restrictedWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

func IsValidChirp(body string) bool {
	return len(body) <= 140
}

func SanitizeChirp(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if restrictedWords[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
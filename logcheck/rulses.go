package logcheck

import (
	"strings"
	"unicode"
)

func IsSmallFrstLetter(literal string) bool {
	letter := []rune(literal)[0]
	return !unicode.IsUpper(letter)
}

func IsOnlyEnglishLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			if !unicode.In(r, unicode.Latin) {
				return false
			}
		}
	}
	return true
}

func nonBannedCharacters(s string, bannedChars []string) bool {
	banned := strings.Join(bannedChars, "")
	for _, r := range s {
		if strings.ContainsRune(banned, r) {
			return false
		}
		if unicode.Is(unicode.S, r) || unicode.Is(unicode.So, r) {
			return false
		}
	}
	return true
}

func nonBannedWords(literal string) bool {
	s := strings.ToLower(literal)
	for _, w := range sensitiveWords {
		if strings.Contains(s, w) {
			return false
		}
	}
	return true
}

var sensitiveWords = []string{
	"password",
	"passwd",
	"token",
	"access_token",
	"refresh_token",
	"api_key",
	"apikey",
	"secret",
	"private_key",
	"ssh_key",
	"cookie",
	"session",
}

var bannedCharacters = []string{
	"!", "@", "#", "$", "%", "^", "&", "*", "(", ")",
	"-", "+", "=", "{", "}", "[", "]", "|", "\\",
	":", ";", "\"", "'", "<", ">", ",", ".", "?", "/",
}

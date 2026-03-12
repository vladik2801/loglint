package logcheck

import (
	"strings"
	"unicode"
)

func sanitizeNoSpecial(msg string, cfg Config) (string, bool) {
	banned := strings.Join(cfg.BannedCharacters, "")

	var out []rune
	changed := false

	for _, r := range msg {
		if unicode.Is(unicode.S, r) || unicode.Is(unicode.So, r) {
			changed = true
			continue
		}
		if strings.ContainsRune(banned, r) {
			changed = true
			continue
		}

		out = append(out, r)
	}

	newMsg := string(out)
	if newMsg != msg {
		changed = true
	}
	return newMsg, changed
}
func buildFixedMessage(msg string, cfg Config) (fixed string, changed bool) {
	fixed = msg

	if cfg.Rules["frstLower"] {
		fixed2, ch := fixLowerFirst(fixed)
		if ch {
			fixed, changed = fixed2, true
		}
	}

	if cfg.Rules["noSpecial"] {
		fixed2, ch := sanitizeNoSpecial(fixed, cfg)
		if ch {
			fixed, changed = fixed2, true
		}
	}

	return fixed, changed
}

func fixLowerFirst(s string) (string, bool) {
	rs := []rune(s)
	if len(rs) == 0 {
		return s, false
	}
	if unicode.IsLetter(rs[0]) && unicode.IsUpper(rs[0]) {
		rs[0] = unicode.ToLower(rs[0])
		return string(rs), true
	}
	return s, false
}

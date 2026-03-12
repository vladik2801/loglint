package logcheck

import (
	"bytes"
	"encoding/json"
	"os"
)

type Config struct {
	Rules            map[string]bool `json:"rules"`
	BannedWords      []string        `json:"banned_words"`
	BannedCharacters []string        `json:"banned_characters"`
}

func loadConfig(pathName string) (Config, error) {
	cfg := DefaultConfig

	data, err := os.ReadFile(pathName)
	if err != nil {
		return cfg, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

var DefaultConfig Config = Config{
	Rules: map[string]bool{
		"frstLower":   true,
		"onlyEng":     true,
		"noSpecial":   true,
		"noSensitive": true,
	},
	BannedWords:      sensitiveWords,
	BannedCharacters: bannedCharacters,
}

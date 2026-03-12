package testpkg

import "log/slog"

func badEng() {
	slog.Info("привет") // want "log message should contain only English letters"
}

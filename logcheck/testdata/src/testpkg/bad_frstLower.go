package testpkg

import "log/slog"

func badLower() {
	slog.Info("Hello world") // want "log message should start with a lowercase letter"
}

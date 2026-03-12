package testpkg

import "log/slog"

func badSensitive() {
	slog.Info("token leaked") // want "log message should not contain sensitive words"
}

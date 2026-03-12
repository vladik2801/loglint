package testpkg

import "log/slog"

func badSpecial() {
	slog.Info("hello!") // want "log message should not contain special characters"
}

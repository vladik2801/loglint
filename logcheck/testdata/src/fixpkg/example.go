package fixpkg

import "log/slog"

func f() {
	slog.Info("Hello!") // want "start with a lowercase letter"
}

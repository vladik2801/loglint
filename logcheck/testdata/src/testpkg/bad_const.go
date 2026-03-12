package testpkg

import "log/slog"

const Msg = "Hello world"

func badConst() {
	slog.Info(Msg) // want "log message should start with a lowercase letter"
}

package logcheck

var supportedLoggers = map[string]bool{
	"log/slog":        true,
	"go.uber.org/zap": true,
}

var LogMethods = map[string]map[string]bool{
	"log/slog": {
		"Debug": true,
		"Info":  true,
		"Warn":  true,
		"Error": true,
	},
	"go.uber.org/zap": {
		"Debug": true,
		"Info":  true,
		"Warn":  true,
		"Error": true,
	},
}

func IsSupportedLogger(pkg string) bool {
	_, ok := supportedLoggers[pkg]
	return ok
}

func IsLoggerMethod(pkg, method string) bool {
	if methods, ok := LogMethods[pkg]; ok {
		return methods[method]
	}
	return false
}

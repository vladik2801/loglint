package logcheck

var configPath string
var enableFix bool

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", "", "path to JSON config")
	Analyzer.Flags.BoolVar(&enableFix, "fix", false, "attach suggested fixes to diagnostics")
}

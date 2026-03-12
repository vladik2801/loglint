package pluginmodule

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/vladik2801/loglint/logcheck"
)

func init() {
	register.Plugin("loglint", New)
}

func New(settings any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

type plugin struct{}

func (*plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		logcheck.Analyzer,
	}, nil
}

func (*plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

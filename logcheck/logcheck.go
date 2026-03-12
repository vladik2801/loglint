package logcheck

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "check for log statements that may be missing context",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	config := DefaultConfig

	if configPath != "" {
		cfg, err := loadConfig(configPath)
		if err != nil {
			return nil, err
		}
		config = cfg
	}
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodefilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodefilter, func(n ast.Node) {
		checkCall(pass, n.(*ast.CallExpr), config)
	})
	return nil, nil
}

func reportRules(pass *analysis.Pass, pos token.Pos, msg string, cfg Config, skipLower, skipSpecial bool) {
	if cfg.Rules["frstLower"] && !skipLower {
		if !IsSmallFrstLetter(msg) {
			pass.Reportf(pos, "loglint: log message should start with a lowercase letter")
		}
	}
	if cfg.Rules["onlyEng"] {
		if !IsOnlyEnglishLetters(msg) {
			pass.Reportf(pos, "loglint: log message should contain only English letters")
		}
	}
	if cfg.Rules["noSpecial"] && !skipSpecial {
		if !nonBannedCharacters(msg, cfg.BannedCharacters) {
			pass.Reportf(pos, "loglint: log message should not contain special characters")
		}
	}
	if cfg.Rules["noSensitive"] {
		if !nonBannedWords(msg) {
			pass.Reportf(pos, "loglint: log message should not contain sensitive words")
		}
	}
}

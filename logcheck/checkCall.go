package logcheck

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

func checkCall(pass *analysis.Pass, call *ast.CallExpr, cfg Config) {
	sell, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	methodName := sell.Sel.Name

	var importPath string

	if ident, ok := sell.X.(*ast.Ident); ok {
		pkgObj, ok := pass.TypesInfo.Uses[ident].(*types.PkgName)
		if !ok {
			return
		}
		importPath = pkgObj.Imported().Path()
	} else if callX, ok := sell.X.(*ast.CallExpr); ok {
		innerSel, ok := callX.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		innerIdent, ok := innerSel.X.(*ast.Ident)
		if !ok {
			return
		}
		pkgObj, ok := pass.TypesInfo.Uses[innerIdent].(*types.PkgName)
		if !ok {
			return
		}
		importPath = pkgObj.Imported().Path()
	} else {
		return
	}
	if !IsSupportedLogger(importPath) {
		return
	}
	if !IsLoggerMethod(importPath, methodName) {
		return
	}
	if len(call.Args) == 0 {
		return
	}

	arg0 := call.Args[0]

	if lit, ok := arg0.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		handleLiteral(pass, lit, cfg)
		return
	}

	handleConst(pass, arg0, cfg)
}

func handleLiteral(pass *analysis.Pass, lit *ast.BasicLit, cfg Config) {
	unquoted, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	lowerBad := cfg.Rules["frstLower"] && !IsSmallFrstLetter(unquoted)
	specialBad := cfg.Rules["noSpecial"] && !nonBannedCharacters(unquoted, cfg.BannedCharacters)

	if enableFix {
		fixed := unquoted

		if cfg.Rules["frstLower"] {
			if s, ch := fixLowerFirst(fixed); ch {
				fixed = s
			}
		}

		if cfg.Rules["noSpecial"] {
			if s, ch := sanitizeNoSpecial(fixed, cfg); ch {
				fixed = s
			}
		}

		if fixed != unquoted && fixed != "" {
			msg := "loglint: log message can be auto-fixed"
			if lowerBad {
				msg = "loglint: log message should start with a lowercase letter"
			} else if specialBad {
				msg = "loglint: log message should not contain special characters"
			}

			pass.Report(analysis.Diagnostic{
				Pos:     lit.Pos(),
				End:     lit.End(),
				Message: msg,
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Apply log message fixes",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     lit.Pos(),
								End:     lit.End(),
								NewText: []byte(strconv.Quote(fixed)),
							},
						},
					},
				},
			})

			reportRules(pass, lit.Pos(), unquoted, cfg, true, true)
			return
		}
	}

	reportRules(pass, lit.Pos(), unquoted, cfg, false, false)
}
func handleConst(pass *analysis.Pass, expr ast.Expr, cfg Config) {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return
	}
	c, ok := pass.TypesInfo.Uses[id].(*types.Const)
	if !ok {
		return
	}
	if c.Val().Kind() != constant.String {
		return
	}

	val := constant.StringVal(c.Val())
	reportRules(pass, id.Pos(), val, cfg, false, false)
}

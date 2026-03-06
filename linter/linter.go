package linters

import (
	"fmt"
	"github.com/golangci/plugin-module-register/register"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
	"strconv"
	"strings"
	"unicode"
)

func init() {
	register.Plugin("awesome_linter", New)
}

type MySettings struct {
	DenyBeginUpper     bool `json:"deny_begin_upper"`
	CheckEnglishOnly   bool `json:"english_only"`
	DenySpecialSymbols bool `json:"deny_special_symbols"`
	CheckSensitive     bool `json:"check_sensitive"`
}

type PluginExample struct {
	settings MySettings
}

func New(settings any) (register.LinterPlugin, error) {

	s, err := register.DecodeSettings[MySettings](settings)
	if err != nil {
		return nil, err
	}

	return &PluginExample{settings: s}, nil
}

func (f *PluginExample) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "awesome_linter",
			Doc:  "some docs",
			Run:  f.run,
		},
	}, nil
}

func (f *PluginExample) GetLoadMode() string {
	// NOTE: the mode can be `register.LoadModeSyntax` or `register.LoadModeTypesInfo`.
	// - `register.LoadModeSyntax`: if the linter doesn't use types information.
	// - `register.LoadModeTypesInfo`: if the linter uses types information.

	return register.LoadModeTypesInfo
}

func (f *PluginExample) run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CallExpr:
				obj := typeutil.Callee(pass.TypesInfo, node)

				if obj == nil || obj.Pkg() == nil {
					return true
				}

				pkgPath := obj.Pkg().Path()

				// I use suffix because of vendor.
				isZap := strings.HasSuffix(pkgPath, "go.uber.org/zap") || strings.HasSuffix(pkgPath, "github.com/uber-go/zap")
				isSlog := pkgPath == "log/slog"

				if !isZap && !isSlog {
					return true
				}

				var lints []analysis.Diagnostic
				switch obj.Name() {
				/**
				func Debug(msg string, args ...any)
				func DebugContext(ctx context.Context, msg string, args ...any)
				func Error(msg string, args ...any)
				func ErrorContext(ctx context.Context, msg string, args ...any)
				func Info(msg string, args ...any)
				func InfoContext(ctx context.Context, msg string, args ...any)
				func Log(ctx context.Context, level Level, msg string, args ...any)
				func LogAttrs(ctx context.Context, level Level, msg string, attrs ...Attr)
				func Warn(msg string, args ...any)
				func WarnContext(ctx context.Context, msg string, args ...any)
				*/
				case "Info", "Debug", "Error", "Warn", "Fatal", "Panic":
					if msg, ok := node.Args[0].(*ast.BasicLit); ok {
						if msg.Kind == token.STRING {
							if f.settings.DenyBeginUpper {
								lints = append(lints, checkLower(node, msg)...)
							}
							if f.settings.CheckEnglishOnly {
								lints = append(lints, checkOnlyEnglish(node, msg)...)
							}
							if f.settings.DenySpecialSymbols {
								lints = append(lints, checkNoSpecialSymbols(node, msg)...)
							}
						}
					}
					if f.settings.CheckSensitive {
						for _, arg := range node.Args[0:] {
							lints = append(lints, checkSensitive(pass, arg)...)
						}
					}
				case "InfoContext", "WarnContext", "ErrorContext", "DebugContext":
					// Zap doesn't have such functions
					if msg, ok := node.Args[1].(*ast.BasicLit); ok {
						if msg.Kind == token.STRING {
							if f.settings.DenyBeginUpper {
								lints = append(lints, checkLower(node, msg)...)
							}
							if f.settings.CheckEnglishOnly {
								lints = append(lints, checkOnlyEnglish(node, msg)...)
							}
							if f.settings.DenySpecialSymbols {
								lints = append(lints, checkNoSpecialSymbols(node, msg)...)
							}
						}
					}
					if f.settings.CheckSensitive {
						for _, arg := range node.Args[1:] {
							lints = append(lints, checkSensitive(pass, arg)...)
						}
					}
				}

				for _, lint := range lints {
					pass.Report(lint)
				}
			}

			return true
		})
	}

	return nil, nil
}

var sensitiveWords = []string{"password", "secret", "api_key", "token", "apikey", "auth_token"}

func isSensitive(name string) bool {
	name = strings.ToLower(name)
	for _, word := range sensitiveWords {
		if strings.Contains(name, word) {
			return true
		}
	}
	return false
}

func checkSensitive(pass *analysis.Pass, node ast.Expr) []analysis.Diagnostic {
	var lints []analysis.Diagnostic
	if bin, ok := node.(*ast.BinaryExpr); ok {
		return append(checkSensitive(pass, bin.X), checkSensitive(pass, bin.Y)...)
	}

	if paren, ok := node.(*ast.ParenExpr); ok {
		return checkSensitive(pass, paren.X)
	}

	if ident, ok := node.(*ast.Ident); ok {
		if isSensitive(ident.Name) {
			lints = append(lints, analysis.Diagnostic{
				Pos:     ident.Pos(),
				End:     ident.End(),
				Message: "log messages should not contain potentially sensitive data.",
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: "do not put sensitive data",
					TextEdits: []analysis.TextEdit{
						{
							ident.Pos(),
							ident.End(),
							[]byte("\"***\""),
						},
					}},
				}})
		}
	}
	return lints
}

func checkLower(node *ast.CallExpr, msg *ast.BasicLit) []analysis.Diagnostic {
	var lints []analysis.Diagnostic
	s, err := strconv.Unquote(msg.Value)
	if err != nil || len(s) == 0 {
		return nil
	}

	runes := []rune(s)
	first := runes[0]

	if unicode.IsUpper(first) {

		runes[0] = unicode.ToLower(first)
		fixedText := string(runes)

		quote := msg.Value[0]
		replacement := fmt.Sprintf("%c%s%c", quote, fixedText, quote)

		lints = append(lints, analysis.Diagnostic{
			Pos:     msg.Pos(),
			End:     msg.End(),
			Message: "log messages must begin with a lowercase letter.",
			SuggestedFixes: []analysis.SuggestedFix{{
				Message: "use lowercase letter on begin",
				TextEdits: []analysis.TextEdit{
					{
						msg.Pos(),
						msg.End(),
						[]byte(replacement),
					},
				}},
			}})
	}
	return lints
}

func checkOnlyEnglish(node *ast.CallExpr, msg *ast.BasicLit) []analysis.Diagnostic {
	var lints []analysis.Diagnostic
	s, err := strconv.Unquote(msg.Value)
	if err != nil {
		return nil
	}
	for _, r := range s {
		if r > unicode.MaxASCII {
			lints = append(lints, analysis.Diagnostic{
				Pos:     msg.Pos(),
				End:     msg.End(),
				Message: "log messages must be in English only.",
			})
			return lints
		}
	}
	return lints
}

func checkNoSpecialSymbols(node *ast.CallExpr, msg *ast.BasicLit) []analysis.Diagnostic {
	var lints []analysis.Diagnostic
	s, err := strconv.Unquote(msg.Value)
	if err != nil {
		return nil
	}
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}

		lints = append(lints, analysis.Diagnostic{
			Pos:     msg.Pos(),
			End:     msg.End(),
			Message: "log messages must not contain special symbols.",
		})
		return lints
	}
	return lints
}

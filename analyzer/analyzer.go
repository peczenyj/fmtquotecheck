package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "fmtquotecheck"
	doc  = name + "verify when safely escape and single quote strings on fmt.Sprintf"
	url  = "https://github.com/peczenyj/fmtquotecheck"
)

func New() *analysis.Analyzer {
	var instance fmtQuoteCheckAnalyzer

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		URL:  url,
		Run:  instance.Run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	return analyzer
}

type fmtQuoteCheckAnalyzer struct{}

func (fa *fmtQuoteCheckAnalyzer) Run(pass *analysis.Pass) (interface{}, error) {
	insp, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(node ast.Node) {
		fa.checkAstNode(pass, node)
	})

	return nil, nil //nolint: nilnil // it is fine return nil,nil here
}

func (fa *fmtQuoteCheckAnalyzer) checkAstNode(pass *analysis.Pass, node ast.Node) {
	if call, ok := node.(*ast.CallExpr); ok {
		fa.checkAstCallExpression(pass, call)
	}
}

func (fa *fmtQuoteCheckAnalyzer) checkAstCallExpression(pass *analysis.Pass,
	call *ast.CallExpr,
) {
	if called, ok := call.Fun.(*ast.SelectorExpr); ok {
		fa.checkAstSelectorExpression(pass, call, called)
	}
}

func (fa *fmtQuoteCheckAnalyzer) checkAstSelectorExpression(pass *analysis.Pass,
	call *ast.CallExpr,
	called *ast.SelectorExpr,
) {
	if expression, ok := called.X.(*ast.Ident); ok {
		fa.checkFullQualifiedFunctionCall(pass, call, called, expression)
	}
}

func (fa *fmtQuoteCheckAnalyzer) checkFullQualifiedFunctionCall(pass *analysis.Pass,
	call *ast.CallExpr,
	called *ast.SelectorExpr,
	expression *ast.Ident,
) {
	fullQualifiedFunctionName := expression.Name + "." + called.Sel.Name

	switch fullQualifiedFunctionName {
	case "fmt.Printf", "fmt.Sprintf", "fmt.Errorf":
		fa.searchForBadQuotedTemplate(pass,
			fullQualifiedFunctionName,
			call.Args[0],
		)
	case "fmt.Fprintf":
		fa.searchForBadQuotedTemplate(pass,
			fullQualifiedFunctionName,
			call.Args[1],
		)
	}
}

func (fa *fmtQuoteCheckAnalyzer) searchForBadQuotedTemplate(pass *analysis.Pass,
	fullQualifiedFunctionName string,
	value ast.Expr,
) {
	if templateLit, ok := value.(*ast.BasicLit); ok && templateLit.Kind == token.STRING {
		template, err := strconv.Unquote(templateLit.Value)
		if err != nil {
			_ = err

			return
		}

		if strings.Contains(template, "'%s'") {
			fix := strings.ReplaceAll(template, "'%s'", "%q")

			msg := "explicit single-quoted '%s' should be replaced by `%q` in " + fullQualifiedFunctionName

			pass.Report(analysis.Diagnostic{
				Pos:     templateLit.Pos(),
				End:     templateLit.End(),
				Message: msg,
				SuggestedFixes: []analysis.SuggestedFix{
					{
						TextEdits: []analysis.TextEdit{
							{
								Pos:     templateLit.Pos(),
								End:     templateLit.End(),
								NewText: []byte(strconv.Quote(fix)),
							},
						},
					},
				},
			})
		}
	}
}

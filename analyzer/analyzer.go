package analyzer

import (
	"errors"
	"flag"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const (
	name = "fmtquotecheck"
	doc  = name + "verify when safely escape and single quote strings on fmt.Sprintf"
	url  = "https://github.com/peczenyj/fmtquotecheck"
)

func New() *analysis.Analyzer {
	var instance fmtQuoteCheckAnalyzer

	instance.SetDefaults()

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		URL:  url,
		Run:  instance.Run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	instance.bindFlags(&analyzer.Flags)

	return analyzer
}

type fmtQuoteCheckAnalyzer struct {
	printfFuncs stringSet
}

func (fa *fmtQuoteCheckAnalyzer) SetDefaults() {
	fa.printfFuncs = stringSet{
		"fmt.Printf":           struct{}{},
		"fmt.Sprintf":          struct{}{},
		"fmt.Errorf":           struct{}{},
		"fmt.Fprintf":          struct{}{},
		"log.Fatalf":           struct{}{},
		"log.Panicf":           struct{}{},
		"log.Printf":           struct{}{},
		"(*log.Logger).Fatalf": struct{}{},
		"(*log.Logger).Panicf": struct{}{},
		"(*log.Logger).Printf": struct{}{},
	}
}

func (fa *fmtQuoteCheckAnalyzer) bindFlags(flagSet *flag.FlagSet) {
	flagSet.Var(&fa.printfFuncs,
		"funcs",
		"full qualified function names to check, split by comma")
}

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
	funcObj, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if !ok {
		return
	}

	fullName := funcObj.Origin().FullName()

	if _, isPrintf := fa.printfFuncs[fullName]; !isPrintf {
		return
	}

	if called, ok := call.Fun.(*ast.SelectorExpr); ok {
		fa.checkAstSelectorExpression(pass, call, called)
	}
}

func (fa *fmtQuoteCheckAnalyzer) checkAstSelectorExpression(pass *analysis.Pass,
	call *ast.CallExpr,
	called *ast.SelectorExpr,
) {
	if expression, ok := called.X.(*ast.Ident); ok {
		fa.checkAstIdentFullQualifiedFunctionCall(pass, call, called, expression)

		return
	}
}

func (fa *fmtQuoteCheckAnalyzer) checkAstIdentFullQualifiedFunctionCall(pass *analysis.Pass,
	call *ast.CallExpr,
	called *ast.SelectorExpr,
	expression *ast.Ident,
) {
	if len(call.Args) <= 1 {
		return
	}

	fullQualifiedFunctionName := expression.Name + "." + called.Sel.Name

	fa.searchForBadQuotedTemplate(pass,
		fullQualifiedFunctionName,
		call.Args[:2],
	)
}

func (fa *fmtQuoteCheckAnalyzer) searchForBadQuotedTemplate(pass *analysis.Pass,
	fullQualifiedFunctionName string,
	values []ast.Expr,
) {
	for _, value := range values {
		if templateLit, ok := value.(*ast.BasicLit); ok && templateLit.Kind == token.STRING {
			fa.checkTemplateLiteral(pass, fullQualifiedFunctionName, templateLit)

			return
		}
	}
}

func (fa *fmtQuoteCheckAnalyzer) checkTemplateLiteral(pass *analysis.Pass,
	fullQualifiedFunctionName string,
	templateLit *ast.BasicLit,
) {
	template, err := strconv.Unquote(templateLit.Value)
	if err != nil {
		_ = err // this should be unreachable

		return
	}

	toSubstitute := strings.Count(template, "'%s'")
	if toSubstitute == 0 {
		return
	}

	msg := "explicit single-quoted '%s' should be replaced by `%q` in "
	msg += fullQualifiedFunctionName

	fix := strconv.Quote(strings.Replace(template, "'%s'", "%q", toSubstitute))

	textEdit := analysis.TextEdit{
		Pos:     templateLit.Pos(),
		End:     templateLit.End(),
		NewText: []byte(fix),
	}

	suggestedFix := analysis.SuggestedFix{
		Message:   "replacing '%s' by `%q`",
		TextEdits: []analysis.TextEdit{textEdit},
	}

	pass.Report(analysis.Diagnostic{
		Pos:     templateLit.Pos(),
		End:     templateLit.End(),
		Message: msg,
		URL:     "https://pkg.go.dev/fmt#pkg-overview",
		SuggestedFixes: []analysis.SuggestedFix{
			suggestedFix,
		},
	})
}

type stringSet map[string]struct{}

var errEmptyString = errors.New("empty string")

func (ss stringSet) Set(value string) error {
	maps.Clear(ss)

	for _, name := range strings.Split(value, ",") {
		if name == "" {
			return errEmptyString
		}

		ss[name] = struct{}{}
	}

	return nil
}

func (ss stringSet) String() string {
	list := make([]string, len(ss))
	cursor := 0

	for name := range ss {
		list[cursor] = name

		cursor++
	}

	sort.Strings(list)

	return strings.Join(list, ",")
}

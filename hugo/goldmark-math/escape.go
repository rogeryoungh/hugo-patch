package math

import (
	"bytes"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type EscapeDollarNode struct {
	ast.BaseInline
}

func (n *EscapeDollarNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindEscapeDollar = ast.NewNodeKind("EscapeDollarNode")

func (n *EscapeDollarNode) Kind() ast.NodeKind {
	return KindEscapeDollar
}

type EscapeDollarParser struct {
}

func (e EscapeDollarParser) Trigger() []byte {
	return []byte{'\\'}
}

func (e EscapeDollarParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, []byte(`\$`)) {
		return nil
	}
	block.Advance(2)
	return &EscapeDollarNode{}
}

type EscapeDollarRenderer struct {
	escape string
}

func (r *EscapeDollarRenderer) renderEscapeDollar(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(r.escape)
		return ast.WalkContinue, nil
	} else {
		return ast.WalkContinue, nil
	}
}

func (r *EscapeDollarRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindEscapeDollar, r.renderEscapeDollar)
}

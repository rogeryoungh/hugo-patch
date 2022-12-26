package math

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type InlineMathNode struct {
	ast.BaseInline
}

type DisplayMathNode struct {
	ast.BaseInline
}

var KindInlineMathNode = ast.NewNodeKind("InlineMathNode")

var KindDisplayMathNode = ast.NewNodeKind("DisplayMathNode")

func (n *InlineMathNode) Kind() ast.NodeKind {
	return KindInlineMathNode
}

func (n *DisplayMathNode) Kind() ast.NodeKind {
	return KindDisplayMathNode
}

func (n *InlineMathNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

func (n *DisplayMathNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

type MathParser struct {
}

func (s *MathParser) Trigger() []byte {
	return []byte{'$'}
}

func (s *MathParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, seg0 := block.PeekLine()
	beg := 0
	for beg < len(line) && line[beg] == '$' {
		beg++
	}
	block.Advance(beg)
	l, pos := block.Position()

	end := 0
	if beg == 1 {
		for {
			c := block.Peek()
			if c == text.EOF {
				block.SetPosition(l, pos)
				return nil
			}
			block.Advance(1)
			end++
			if c == '$' {
				break
			}
		}
		node := &InlineMathNode{}
		seg := text.NewSegment(seg0.Start+beg, seg0.Start+end)
		node.AppendChild(node, ast.NewRawTextSegment(seg))
		return node
	} else if beg == 2 {
		c := byte(0)
		count := 0
		for {
			c = block.Peek()
			if c == text.EOF {
				block.SetPosition(l, pos)
				return nil
			}
			block.Advance(1)
			end++
			if c == '$' {
				if count == 0 {
					count = 1
					continue
				}
			}
			if count == 1 {
				if c == '$' {
					break
				} else {
					continue
				}
			}
		}
		node := &DisplayMathNode{}
		seg := text.NewSegment(seg0.Start+beg, seg0.Start+end)
		node.AppendChild(node, ast.NewRawTextSegment(seg))
		return node
	}
	return nil
}

type InlineMathRenderer struct {
	startDelim string
	endDelim   string
}

type DisplayMathRenderer struct {
	startDelim string
	endDelim   string
}

func (r *InlineMathRenderer) renderInlineMath(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(r.startDelim)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			w.Write(segment.Value(source))
		}
		return ast.WalkSkipChildren, nil
	} else {
		w.WriteString(r.endDelim)
		return ast.WalkContinue, nil
	}
}

func (r *DisplayMathRenderer) renderDisplayMath(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(r.startDelim)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			w.Write(segment.Value(source))
		}
		return ast.WalkSkipChildren, nil
	} else {
		w.WriteString(r.endDelim)
		return ast.WalkContinue, nil
	}
}

func (r *InlineMathRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindInlineMathNode, r.renderInlineMath)
}

func (r *DisplayMathRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindDisplayMathNode, r.renderDisplayMath)
}

package math

import (
	"bytes"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)


type InlineMath struct {
	ast.BaseInline
}

func (n *InlineMath) Inline() {}

func (n *InlineMath) IsBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*ast.Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

func (n *InlineMath) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindInlineMath = ast.NewNodeKind("InlineMath")

func (n *InlineMath) Kind() ast.NodeKind {
	return KindInlineMath
}

func NewInlineMath() *InlineMath {
	return &InlineMath{
		BaseInline: ast.BaseInline{},
	}
}


type inlineMathParser struct {
}

var defaultInlineMathParser = &inlineMathParser{}

func NewInlineMathParser() parser.InlineParser {
	return defaultInlineMathParser
}

func (s *inlineMathParser) Trigger() []byte {
	return []byte{'$'}
}

func (s *inlineMathParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, startSegment := block.PeekLine()
	opener := 0
	for ; opener < len(line) && line[opener] == '$'; opener++ {
	}
	block.Advance(opener)
	l, pos := block.Position()
	node := NewInlineMath()
	for {
		line, segment := block.PeekLine()
		if line == nil {
			block.SetPosition(l, pos)
			return ast.NewTextSegment(startSegment.WithStop(startSegment.Start + opener))
		}
		for i := 0; i < len(line); i++ {
			c := line[i]
			if c == '$' {
				oldi := i
				for ; i < len(line) && line[i] == '$'; i++ {
				}
				closure := i - oldi
				if closure == opener && (i+1 >= len(line) || line[i+1] != '$') {
					segment := segment.WithStop(segment.Start + i - closure)
					if !segment.IsEmpty() {
						node.AppendChild(node, ast.NewRawTextSegment(segment))
					}
					block.Advance(i)
					goto end
				}
			}
		}
		if !util.IsBlank(line) {
			node.AppendChild(node, ast.NewRawTextSegment(segment))
		}
		block.AdvanceLine()
	}
end:

	if !node.IsBlank(block.Source()) {
		// trim first halfspace and last halfspace
		segment := node.FirstChild().(*ast.Text).Segment
		shouldTrimmed := true
		if !(!segment.IsEmpty() && block.Source()[segment.Start] == ' ') {
			shouldTrimmed = false
		}
		segment = node.LastChild().(*ast.Text).Segment
		if !(!segment.IsEmpty() && block.Source()[segment.Stop-1] == ' ') {
			shouldTrimmed = false
		}
		if shouldTrimmed {
			t := node.FirstChild().(*ast.Text)
			segment := t.Segment
			t.Segment = segment.WithStart(segment.Start + 1)
			t = node.LastChild().(*ast.Text)
			segment = node.LastChild().(*ast.Text).Segment
			t.Segment = segment.WithStop(segment.Stop - 1)
		}

	}
	return node
}

type InlineMathRenderer struct {
	startDelim string
	endDelim   string
}

func (r *InlineMathRenderer) renderInlineMath(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(`` + r.startDelim)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := segment.Value(source)
			if bytes.HasSuffix(value, []byte("\n")) {
				w.Write(value[:len(value)-1])
				if c != n.LastChild() {
					w.Write([]byte(" "))
				}
			} else {
				w.Write(value)
			}
		}
		return ast.WalkSkipChildren, nil
	}
	_, _ = w.WriteString(r.endDelim + ``)
	return ast.WalkContinue, nil
}

func (r *InlineMathRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindInlineMath, r.renderInlineMath)
}


func NewInlineMathRenderer(start, end string) renderer.NodeRenderer {
	return &InlineMathRenderer{start, end}
}

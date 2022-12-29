package math

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type MathBlock struct {
	ast.BaseBlock
}

var KindMathBlock = ast.NewNodeKind("MathBLock")

func NewMathBlock() *MathBlock {
	return &MathBlock{}
}

func (n *MathBlock) Dump(source []byte, level int) {
	m := map[string]string{}
	ast.DumpHelper(n, source, level, m, nil)
}

func (n *MathBlock) Kind() ast.NodeKind {
	return KindMathBlock
}

func (n *MathBlock) IsRaw() bool {
	return true
}

type mathJaxBlockParser struct {
}

var defaultMathJaxBlockParser = &mathJaxBlockParser{}

type mathBlockData struct {
	indent int
}

var mathBlockInfoKey = parser.NewContextKey()

func NewMathJaxBlockParser() parser.BlockParser {
	return defaultMathJaxBlockParser
}

func (b *mathJaxBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos == -1 {
		return nil, parser.NoChildren
	}
	if line[pos] != '$' {
		return nil, parser.NoChildren
	}
	i := pos
	for ; i < len(line) && line[i] == '$'; i++ {
	}
	if i-pos < 2 {
		return nil, parser.NoChildren
	}
	pc.Set(mathBlockInfoKey, &mathBlockData{indent: pos})
	node := NewMathBlock()
	return node, parser.NoChildren
}

func (b *mathJaxBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	data := pc.Get(mathBlockInfoKey).(*mathBlockData)
	w, pos := util.IndentWidth(line, 0)
	if w < 4 {
		i := pos
		for ; i < len(line) && line[i] == '$'; i++ {
		}
		length := i - pos
		if length >= 2 && util.IsBlank(line[i:]) {
			reader.Advance(segment.Stop - segment.Start - segment.Padding)
			return parser.Close
		}
	}

	pos, padding := util.IndentPositionPadding(line, 0, 0, data.indent)
	if pos < 0 {
		pos = util.FirstNonSpacePosition(line)
	}
	if padding < 0 {
		padding = 0
	}
	seg := text.NewSegmentPadding(segment.Start+pos, segment.Stop, padding)
	node.Lines().Append(seg)
	reader.AdvanceAndSetPadding(segment.Stop-segment.Start-pos-1, padding)
	return parser.Continue | parser.NoChildren
}

func (b *mathJaxBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	pc.Set(mathBlockInfoKey, nil)
}

func (b *mathJaxBlockParser) CanInterruptParagraph() bool {
	return true
}

func (b *mathJaxBlockParser) CanAcceptIndentedLine() bool {
	return false
}

func (b *mathJaxBlockParser) Trigger() []byte {
	return nil
}

type MathBlockRenderer struct {
	startDelim string
	endDelim   string
}

func NewMathBlockRenderer(start, end string) renderer.NodeRenderer {
	return &MathBlockRenderer{start, end}
}

func (r *MathBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMathBlock, r.renderMathBlock)
}

func (r *MathBlockRenderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		w.Write(line.Value(source))
	}
}

func (r *MathBlockRenderer) renderMathBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*MathBlock)
	if entering {
		_, _ = w.WriteString("<p>" + r.startDelim + "\n")

		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString(r.endDelim + "</p>\n")
	}
	return ast.WalkContinue, nil
}

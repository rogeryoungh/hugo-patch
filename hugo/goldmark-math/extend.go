package math

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type latex struct {
	inlineStartDelim string
	inlineEndDelim   string
	blockStartDelim  string
	blockEndDelim    string
}

type Option interface {
	SetOption(e *latex)
}

type withInlineDelim struct {
	start string
	end   string
}

type withBlockDelim struct {
	start string
	end   string
}

func WithInlineDelim(start string, end string) Option {
	return &withInlineDelim{start, end}
}

func (o *withInlineDelim) SetOption(e *latex) {
	e.inlineStartDelim = o.start
	e.inlineEndDelim = o.end
}

func WithBlockDelim(start string, end string) Option {
	return &withBlockDelim{start, end}
}

func (o *withBlockDelim) SetOption(e *latex) {
	e.blockStartDelim = o.start
	e.blockEndDelim = o.end
}

var LaTeX = &latex{
	inlineStartDelim: `$`,
	inlineEndDelim:   `$`,
	blockStartDelim:  `\[`,
	blockEndDelim:    `\]`,
}

func NewLaTeX(opts ...Option) *latex {
	r := &latex{
		inlineStartDelim: `$`,
		inlineEndDelim:   `$`,
		blockStartDelim:  `\[`,
		blockEndDelim:    `\]`,
	}

	for _, o := range opts {
		o.SetOption(r)
	}
	return r
}

func (e *latex) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&EscapeDollarParser{}, 501),
		util.Prioritized(&MathParser{}, 502),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&EscapeDollarRenderer{"&dollar;"}, 500),
		util.Prioritized(&InlineMathRenderer{"$", "$"}, 501),
		util.Prioritized(&DisplayMathRenderer{"\\[", "\\]"}, 502),
	))
}

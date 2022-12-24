package math

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yuin/goldmark"

	"github.com/stretchr/testify/assert"
)

type latexTestCase struct {
	d   string // test description
	in  string // input markdown source
	out string // expected output html
}

func TestMath_One(t *testing.T) {
	s := []byte("$")
	print(s[0])
	out, err := renderMarkdown([]byte(" $$\n1+2\n$$"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "<p>\\[\n1+2\n\\]</p>", strings.TrimSpace(string(out)))
}

func TestMath(t *testing.T) {

	tests := []latexTestCase{
		{
			d:   "plain text",
			in:  "foo",
			out: `<p>foo</p>`,
		},
		{
			d:   "bold",
			in:  "**foo**",
			out: `<p><strong>foo</strong></p>`,
		},
		{
			d:   "latex inline",
			in:  "$1+2$",
			out: `<p>$1+2$</p>`,
		},
		{
			d:  "latex display",
			in: "$$\n1+2\n$$",
			out: "<p>\\[\n1+2\n\\]</p>",
		},
		{
			// this input previously triggered a panic in block.go
			d:   "list-begin",
			in:  "*foo\n  ",
			out: "<p>*foo</p>",
		},
		{
			d:   "latex in em",
			in:  "_x_\n\n$_x_$\n\n$x_xx_x$\n\n $\\bf{f}_{1} $\\bf{f}_{1}",
			out: "<p><em>x</em></p>\n<p>$_x_$</p>\n<p>$x_xx_x$</p>\n<p>$\\bf{f}_{1} $\\bf{f}_{1}</p>",
		},
		{
			d:   "latex in code",
			in:  "`_x_`  `$_x_$`\n",
			out: "<p><code>_x_</code>  <code>$_x_$</code></p>",
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d: %s", i, tc.d), func(t *testing.T) {
			out, err := renderMarkdown([]byte(tc.in))
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.out, strings.TrimSpace(string(out)))
		})
	}

}

func renderMarkdown(src []byte) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(LaTeX),
	)

	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

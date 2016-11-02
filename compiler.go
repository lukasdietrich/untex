package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

var (
	texEscaper = strings.NewReplacer(
		`\`, `\textbackslash{}`,
		`~`, `\textasciitilde{}`,
		`^`, `\textasciicircum{}`,
		`&`, `\&`,
		`%`, `\%`,
		`$`, `\$`,
		`#`, `\#`,
		`_`, `\_`,
		`{`, `\{`,
		`}`, `\}`,
	)

	requiredPackages = []Package{
		{Name: "tabularx"},
		{Name: "color"},
		{Name: "enumerate"},
		{Name: "hyperref"},
	}
)

type Meta map[string]string

type Compiler struct {
	w  io.Writer
	fs *Filesystem

	dir  string
	seen map[string]bool

	meta Meta
	tmpl *Template
	errc chan<- error

	root  bool
	nl    int
	tabs  int
	block StringStack
}

func (c *Compiler) Compile(src string) error {
	if c.root {
		abs, err := filepath.Abs(src)
		if err != nil {
			return err
		}
		src = abs
	} else if c.seen[src] {
		return fmt.Errorf("multiple imports of '%s' detected", src)
	}

	buf, err := c.fs.ReadFile(src)
	if err != nil {
		return err
	}

	if c.root {
		c.meta = make(Meta)
		c.seen = make(map[string]bool)
	}

	c.dir, c.seen[src] = filepath.Dir(src), true

	g := Grammar{
		Buffer:   string(buf),
		Compiler: c,
	}
	g.Init()
	err = g.Parse()
	if err != nil {
		return err
	}

	errc := make(chan error)
	c.errc = errc

	go func() {
		g.Execute()
		close(errc)
	}()

	return <-errc
}

func (c *Compiler) Import(src string) {
	compiler := Compiler{
		w:    c.w,
		fs:   c.fs,
		seen: c.seen,
	}
	src = filepath.Clean(filepath.Join(c.dir, src))

	if err := compiler.Compile(src); err != nil {
		c.errc <- err
	}
}

func (c *Compiler) Begin() {
	if !c.root {
		return
	}

	tmpl, err := GetTemplate(c.meta["template"])
	if err != nil {
		c.errc <- err
		return
	}

	for _, p := range requiredPackages {
		tmpl.AddPackage(p)
	}

	if err = tmpl.ApplyMeta(c.meta); err != nil {
		c.errc <- err
		return
	}

	if err = tmpl.WritePrefix(c.w); err != nil {
		c.errc <- err
		return
	}

	c.tmpl = tmpl
}

func (c *Compiler) End() {
	if !c.root {
		return
	}

	if err := c.tmpl.WriteSuffix(c.w); err != nil {
		c.errc <- err
		return
	}
}

func (c *Compiler) SetMeta(k, v string) {
	c.meta[k] = v
}

func (c *Compiler) NewLine() {
	if c.nl > 1 {
		return
	}

	c.AddLatex("\n")
	c.nl++
}

func (c *Compiler) BeginBlock(n string) {
	c.AddLatex(`\begin{` + n + "}\n")
	c.block.Push(n)
}

func (c *Compiler) EndBlock() {
	c.AddLatex(`\end{` + c.block.Pop() + "}\n")
}

func (c *Compiler) EndAllBlocks() {
	for !c.block.Empty() {
		c.EndBlock()
	}
}

func (c *Compiler) AddLatex(t string) {
	c.nl = 0
	_, err := fmt.Fprint(c.w, t)
	if err != nil {
		c.errc <- err
	}
}

func (c *Compiler) AddText(t string) {
	c.AddLatex(texEscaper.Replace(t))
}

func (c *Compiler) AddSection(t string, l int) {
	c.AddLatex(fmt.Sprintf(
		`\%ssection{%s}`,
		strings.Repeat("sub", l),
		t,
	))
}

func (c *Compiler) AddEmph(t string) {
	c.AddLatex(`\textit{`)
	c.AddText(t)
	c.AddLatex(`}`)
}

func (c *Compiler) AddBold(t string) {
	c.AddLatex(`\textbf{`)
	c.AddText(t)
	c.AddLatex(`}`)
}

func (c *Compiler) AddLink(url, text string) {
	c.AddLatex(`\href{`)
	c.AddText(url)
	c.AddLatex(`}{`)
	c.AddText(text)
	c.AddLatex(`}`)
}

func (c *Compiler) AddListItem(list, label string) {
	d := c.tabs + 1 - c.block.Size()
	if d > 0 {
		for ; d > 0; d-- {
			c.BeginBlock(list)
		}
	} else {
		for ; d < 0; d++ {
			c.EndBlock()
		}
	}

	if c.block.Peek() != list {
		c.EndBlock()
		c.BeginBlock(list)
	}

	if label != "" {
		c.AddLatex(`\item[`)
		c.AddText(label)
		c.AddLatex(`] `)
	} else {
		c.AddLatex(`\item `)
	}
}

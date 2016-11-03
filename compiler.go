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
	w     io.Writer
	fs    *Filesystem
	trace StringStack

	meta Meta
	tmpl *Template
	errc chan error

	nl    int
	tabs  int
	block StringStack
}

func (c *Compiler) compile(src string) error {
	abs, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	if c.trace.Contains(abs) {
		return fmt.Errorf("recursive import of '%s'", abs)
	}

	buf, err := c.fs.ReadFile(abs)
	if err != nil {
		return err
	}

	c.trace.Push(abs)
	defer c.trace.Pop()

	g := Grammar{
		Buffer:   string(buf),
		Compiler: c,
	}
	g.Init()
	if err := g.Parse(); err != nil {
		return err
	}

	g.Execute()
	return nil
}

func (c *Compiler) Compile(src string) error {
	c.meta = make(Meta)
	c.errc = make(chan error)

	go func() {
		c.compile(src)
		close(c.errc)
	}()

	return <-c.errc
}

func (c *Compiler) Import(src string) {
	err := c.compile(filepath.Join(filepath.Dir(c.trace.Peek()), src))
	if err != nil {
		c.errc <- err
	}
}

func (c *Compiler) Begin() {
	if c.trace.Size() > 1 {
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
	if c.trace.Size() > 1 {
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

	c.Write("\n")
	c.nl++
}

func (c *Compiler) BeginBlock(n string) {
	c.Write(`\begin{` + n + "}\n")
	c.block.Push(n)
}

func (c *Compiler) EndBlock() {
	c.Write(`\end{` + c.block.Pop() + "}\n")
}

func (c *Compiler) EndAllBlocks() {
	for !c.block.Empty() {
		c.EndBlock()
	}
}

func (c *Compiler) Write(t string) {
	c.nl = 0
	_, err := fmt.Fprint(c.w, t)
	if err != nil {
		c.errc <- err
	}
}

func (c *Compiler) AddText(t string) {
	c.Write(texEscaper.Replace(t))
}

func (c *Compiler) AddLatex(t string) {
	c.Write(strings.Trim(t, "\n"))
}

func (c *Compiler) AddSection(t string, l int) {
	c.Write(fmt.Sprintf(
		`\%ssection{%s}`,
		strings.Repeat("sub", l),
		t,
	))
}

func (c *Compiler) AddEmph(t string) {
	c.Write(`\textit{`)
	c.AddText(t)
	c.Write(`}`)
}

func (c *Compiler) AddBold(t string) {
	c.Write(`\textbf{`)
	c.AddText(t)
	c.Write(`}`)
}

func (c *Compiler) AddLink(url, text string) {
	c.Write(`\href{`)
	c.AddText(url)
	c.Write(`}{`)
	c.AddText(text)
	c.Write(`}`)
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
		c.Write(`\item[`)
		c.AddText(label)
		c.Write(`] `)
	} else {
		c.Write(`\item `)
	}
}

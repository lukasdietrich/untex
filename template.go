package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var (
	prefix, suffix *template.Template

	templateBuffer bytes.Buffer
	whitespace     = regexp.MustCompile(`^(\s*)`)
)

func init() {
	load := func(name string) *template.Template {
		f, _ := assets.Open(name)
		defer f.Close()
		b, _ := ioutil.ReadAll(f)
		return template.Must(template.New(name).Parse(string(b)))
	}

	prefix = load("prefix.tmpl")
	suffix = load("suffix.tmpl")
}

type Document struct {
	Type     string `xml:"type,attr"`
	Options  string `xml:"options,attr"`
	Prefix   string `xml:"prefix"`
	Suffix   string `xml:"suffix"`
	Preamble string `xml:"preamble"`
}

type Package struct {
	Name    string `xml:"name,attr"`
	Options string `xml:"options,attr"`
}

type Asset struct {
	Path   string `xml:"path,attr"`
	Base64 []byte `xml:",innerxml"`
}

type Template struct {
	XMLName  xml.Name  `xml:"template"`
	Document Document  `xml:"document"`
	Packages []Package `xml:"packages>package"`
	Assets   []Asset   `xml:"assets>asset"`
}

func ReadTemplate(r io.Reader) (*Template, error) {
	var t Template
	return &t, xml.NewDecoder(r).Decode(&t)
}

func GetTemplate(src string) (*Template, error) {
	var (
		r   io.ReadCloser
		err error
	)

	switch {
	case src == "":
		r, err = assets.Open("default.xml")
	case strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://"):
		res, _err := http.Get(src)
		r, err = res.Body, _err
	default:
		r, err = os.Open(src)
	}

	if err != nil {
		return nil, err
	}

	defer r.Close()
	return ReadTemplate(r)
}

func (t *Template) AddPackage(p Package) {
	for _, k := range t.Packages {
		if k.Name == p.Name {
			return
		}
	}

	t.Packages = append(t.Packages, p)
}

func (t *Template) ApplyMeta(m Meta) error {
	t.Document.Preamble = trimIndentation(t.Document.Preamble)
	t.Document.Prefix = trimIndentation(t.Document.Prefix)
	t.Document.Suffix = trimIndentation(t.Document.Suffix)

	var (
		err error
		s   = []*string{
			&t.Document.Preamble,
			&t.Document.Prefix,
			&t.Document.Suffix,
			&t.Document.Options,
			&t.Document.Type,
		}
	)

	for _, v := range s {
		*v, err = execTextTemplate(*v, m)
		if err != nil {
			return err
		}
	}

	for i, p := range t.Packages {
		p.Name, err = execTextTemplate(p.Name, m)
		if err != nil {
			return err
		}

		p.Options, err = execTextTemplate(p.Options, m)
		if err != nil {
			return err
		}

		t.Packages[i] = p
	}

	return nil
}

func (t *Template) WritePrefix(w io.Writer) error {
	return prefix.Execute(w, t)
}

func (t *Template) WriteSuffix(w io.Writer) error {
	return suffix.Execute(w, t)
}

func execTextTemplate(tmpl string, values map[string]string) (string, error) {
	t, err := template.New("inline").Parse(tmpl)
	if err != nil {
		return "", err
	}

	templateBuffer.Reset()
	err = t.Execute(&templateBuffer, values)
	if err != nil {
		return "", err
	}

	return templateBuffer.String(), nil
}

func trimIndentation(text string) string {
	if text == "" {
		return text
	}

	lines, ind := strings.Split(text, "\n"), ""

	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			pind := whitespace.FindString(l)
			if ind == "" || len(pind) < len(ind) {
				ind = pind
			}
		}
	}

	for i, l := range lines {
		if strings.HasPrefix(l, ind) {
			lines[i] = l[len(ind):]
		}
	}

	return strings.Join(lines, "\n")
}

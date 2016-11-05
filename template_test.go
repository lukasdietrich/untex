package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateExact(t *testing.T) {
	tmpl := Template{
		Document: Document{
			Type:     "{{.DocType}}",
			Options:  "{{.DocOptions}}",
			Prefix:   "prefix",
			Suffix:   "suffix",
			Preamble: "preamble",
		},
		Packages: []Package{
			Package{
				Name:    "{{.PackageName}}",
				Options: "{{.PackageOptions}}",
			},
		},
	}

	tmpl.ApplyMeta(Meta{
		"DocType":        "test-type",
		"DocOptions":     "test-options",
		"PackageName":    "package1",
		"PackageOptions": "fancy-mode",
	})

	var buf bytes.Buffer
	tmpl.WritePrefix(&buf)
	tmpl.WriteSuffix(&buf)

	expected := []string{
		`\documentclass[test-options]{test-type}`,
		``,
		`\usepackage[fancy-mode]{package1}`,
		`preamble`,
		`\begin{document}`,
		`prefix`,
		`suffix`,
		`\end{document}`,
		``,
	}

	assert.Equal(t, strings.Join(expected, "\n"), buf.String(),
		"generated output should match")
}

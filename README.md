# Untex

A *clean* Latex preprocessor.

## 1. Build Requirements

- GO `1.7+`
- PEG generator has to be in the `PATH` (<https://github.com/pointlander/peg>)

```sh
go get github.com/pointlander/peg
go get github.com/shurcooL/vfsgen

go generate

go get
go build
```

## 2. Documentation

### 2.1 Templates

Templates can link to local paths or `http(s)` urls.
If no template is set, a default one is provided from `assets/default.xml`.

The `preamble`, `prefix`, `suffix` and every attribute is rendered as a 
[text/template](https://golang.org/pkg/text/template/) with every value of the
meta section available.

```xml
<template>
    <!-- translates to \documentclass[options]{type} -->
    <document type="scrartcl" options="a4paper,10pt">
        <preamble>
            inserted at the end of the preamble
        </preamble>
        <prefix>
            inserted right after \begin{document}
        </prefix>
        <suffix>
            inserted right before \end{document}
        </suffix>
    </document>
    <packages>
        <!-- translates to \usepackage[options]{name} -->
        <package name="inputenc" options="utf8" />
    </packages>
</template>

```

### 2.2 Syntax

The root file can (and should) have a meta section.
This section contains of a key-value configuration, that will be used 
to render the template.

```
---
template:   mytemplate.xml

title:      My first Untex document
---
```

#### 2.2.1 Imports

Subfiles can be imported and trandformed as well.

```
@import(sections/my-other-file.tex)
```

#### 2.2.2 Headlines

Headlines start with one to three `#`'s and get translated to 
`\section`, `\subsection` or `\subsubsection` respectively.

```markdown
# This is a section

## This is a subsection
```

#### 2.2.3 Lists

There are three types of lists, that can be nested into each other.

```
- This is an unordered list
    1. This is an
    2. ordered list

Red) This is a named list
    Blue) :-)
```

#### 2.2.4 Plain LaTeX

Plain LaTeX can be used as is, if needed.

```
%%%

\begin{tabular}{c|c|c}
1 & 2 & 3
\end{tabular}

%%%
```

#### 2.2.5 Inlines

In paragraphs, as well as in list items the following can be used.

```
*bold text*
/italic text/
$ inline math $
$$ display math $$
[link text](url)
```

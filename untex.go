package main

//go:generate peg -switch grammar.peg
//go:generate go run assets.go

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"

	"github.com/blang/vfs"
)

const (
	STDOUT = "stdout"
)

func main() {
	var (
		in     string
		out    string
		writer io.Writer

		fs = NewFilesystem(vfs.OS())
	)

	flag.StringVar(&out, "output", STDOUT, "output file")
	flag.Parse()
	in = flag.Arg(0)

	if in == "" {
		log.Fatal("you need to specify a source file")
	}

	if out == STDOUT {
		writer = os.Stdout
		log.SetOutput(os.Stderr)
	} else {
		f, err := fs.Create(out)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		writer = f
	}

	var (
		buffer   = bufio.NewWriter(writer)
		compiler = Compiler{
			w:  buffer,
			fs: fs,
		}
	)

	if err := compiler.Compile(in); err != nil {
		log.Fatal(err)
	}

	if err := buffer.Flush(); err != nil {
		log.Fatal(err)
	}
}

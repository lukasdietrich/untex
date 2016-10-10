package main

//go:generate peg -switch grammar.peg
//go:generate go run assets.go

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
)

const (
	STDOUT = "stdout"
)

func main() {
	var (
		in     string
		out    string
		writer io.Writer
		buffer *bufio.Writer
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
		f, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		writer = f
	}

	buffer = bufio.NewWriter(writer)

	err := (&Compiler{w: buffer, root: true}).Compile(in)
	if err != nil {
		log.Fatal(err)
	}

	err = buffer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

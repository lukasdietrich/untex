// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir("assets/"), vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}
}

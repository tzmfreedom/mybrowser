package main

import (
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/mybrowser/parser"
)

func main() {
	src, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	p := parser.NewHtmlParser(string(src))
	node, err := p.Parse()
	if err != nil {
		panic(err)
	}
	pp.Println(node)
}

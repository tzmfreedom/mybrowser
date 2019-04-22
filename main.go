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
	htmlParser := parser.NewHtmlParser(string(src))
	dom, err := htmlParser.Parse()
	if err != nil {
		panic(err)
	}
	cssParser := parser.NewCssParser(`
div{
	display: block;
	color: red;
}`)
	css, err := cssParser.Parse()
	node := parser.NewStyledNode(dom, css)
	pp.Println(node)
}

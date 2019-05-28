package main

import (
	"github.com/IvoryRaptor/iotbox/iotql"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type j2xConvert struct {
	*iotql.BaseiotqlListener
	xml map[antlr.Tree]string
}

func NewJ2xConvert() *j2xConvert {
	return &j2xConvert{
		&iotql.BaseiotqlListener{},
		make(map[antlr.Tree]string),
	}
}

func main() {
	f, err := os.Open("iotql/aa.sql")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	// Setup the input
	is := antlr.NewInputStream(string(content))

	// Create lexter
	lexer := iotql.NewiotqlLexer(is)

	stream := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)

	// Create parser and tree
	//p := json2xml.NewJSONParser(stream)
	p := iotql.NewiotqlParser(stream)
	p.BuildParseTrees = true

	println(p.Name().GetText())
	var tree = p.Sql_stmt()
	j2x := NewJ2xConvert()
	antlr.ParseTreeWalkerDefault.Walk(j2x, tree)

	log.Println(j2x)
}

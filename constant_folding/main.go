package main

import (
	"flag"
	"fmt"
	"go/ast"
	"log"

	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Usage: <prog> path")
	}

	path := flag.Arg(0)
	conf := &packages.Config{
		Mode: ^0,
		Dir:  path,
	}
	pkgs, err := packages.Load(conf, path)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		walk(pkg)
	}
}

func walk(pkg *packages.Package) {
	for _, file := range pkg.Syntax {
		walkFile(pkg, file)
	}
}

func walkFile(pkg *packages.Package, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch ident := n.(type) {
		case *ast.Ident:
			fmt.Println(ident)
			val := pkg.TypesInfo.Types[ident].Value
			fmt.Println(val)

		case *ast.BasicLit:
			fmt.Println(ident.Value)
		}
		return true
	})
}

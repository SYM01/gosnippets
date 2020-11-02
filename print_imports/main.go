package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

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
		fmt.Printf("@import %s\n", pkg.Name)
		showPkgs(pkg, 3, 0)
	}
}

func showPkgs(pkg *packages.Package, maxDepths, ident int) {
	if maxDepths < 0 {
		return
	}

	// print files in current package
	for _, syn := range pkg.Syntax {
		fmt.Printf("%s - %s\n", strings.Repeat("\t", ident), pkg.Fset.Position(syn.Pos()))
	}

	// print imported packages
	for name, dep := range pkg.Imports {
		fmt.Printf("%s @import %s\n", strings.Repeat("\t", ident), name)
		showPkgs(dep, maxDepths-1, ident+1)
	}
}

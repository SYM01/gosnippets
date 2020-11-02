package main

import (
	"flag"
	"fmt"
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
		showPkgs(pkg, 3, 0)
	}
}

func showPkgs(pkg *packages.Package, maxDepths, ident int) {
	if maxDepths < 0 {
		return
	}

	// print module info for current pacckage
	if mod := pkg.Module; mod != nil {
		fmt.Printf("@package %s\n", pkg.ID)
		showModuel(mod)
	}

	// trace imported packages
	for _, dep := range pkg.Imports {
		showPkgs(dep, maxDepths-1, ident+1)
	}
}

func showModuel(mod *packages.Module) {
	fmt.Printf("# %s\n", mod.Path)
	fmt.Printf("\tVersion: %s\n", mod.Version)
	if mod.Replace != nil {
		fmt.Printf("\tReplace: %s @%s\n", mod.Replace.Path, mod.Replace.Version)
	}
	fmt.Printf("\tMain: %v\n", mod.Main)
	fmt.Printf("\tDir: %s\n", mod.Dir)
	fmt.Printf("\tGoMod: %s\n", mod.GoMod)
	fmt.Printf("\tGoVersion: %s\n", mod.GoVersion)
	fmt.Printf("\tError: %v\n", mod.Error)
}

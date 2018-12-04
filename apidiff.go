package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/tools/go/packages"

	"golang.org/x/exp/apidiff"
)

func main() {
	os.Exit(main1())
}

func main1() int {
	if err := mainerr(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func mainerr() error {
	flag.Parse()

	if len(flag.Args()) != 2 {
		return fmt.Errorf("need exactly two arguments")
	}

	bef, err := loadSinglePkgFromDir(flag.Args()[0])
	if err != nil {
		return err
	}
	aft, err := loadSinglePkgFromDir(flag.Args()[1])
	if err != nil {
		return err
	}

	changes := apidiff.Changes(bef.Types, aft.Types)

	fmt.Printf("%v", changes)

	return nil
}

func loadSinglePkgFromDir(dir string) (*packages.Package, error) {
	conf := &packages.Config{
		Dir:  dir,
		Mode: packages.LoadTypes,
	}

	pkgs, err := packages.Load(conf, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to go/packages.Load (in %v): %v", dir, err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("dir %v resolved to multiple (test) go/packages.Package", dir)
	}

	return pkgs[0], nil
}

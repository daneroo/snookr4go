package main

import (
	"flag"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// see https://lawlessguy.wordpress.com/2013/07/23/filling-a-slice-using-command-line-flags-in-go-golang/ to use multiple arguments.
// or use this alternative package
// var root = flag.Bool("root", '', "a root to traverse")

func main() {
	flag.Parse()
	roots := flag.Args()

	exif.RegisterParsers(mknote.All...)

	// root := flag.Arg(0)
	if len(roots) == 0 {
		fmt.Printf("Root folder(s) not specified.")
		os.Exit(1)
	}

	for _, root := range roots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			fmt.Printf("Root folder not found: %s\n", root)
			continue
		}
		err := filepath.Walk(root, visit)
		if err != nil {
			fmt.Printf("Walk(%s) returned %v\n", root, err)
		} else {
			fmt.Printf("Walk(%s) is done\n", root)

		}
	}

}

func visit(filename string, f os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	// fmt.Printf("f: %v\n", f)
	// fmt.Printf("Visited: %s\n", filename)
	if !f.IsDir() && strings.ToLower(path.Ext(filename)) == ".jpg" {
		exifOne(filename)
	}
	return nil
}

func exifOne(name string) error {
	f, err := os.Open(name)
	if err != nil {
		log.Printf("err on %v: %v", name, err)
		return err
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Printf("err on %v: %v", name, err)
		return err
	}

	fmt.Printf("  ---- Image '%v' ----\n", name)
	stamp, err := x.DateTime() // normally, don't ignore errors!
	if err != nil {
		fmt.Printf("  Date: %v\n", stamp)
	}
	camModel, err := x.Get(exif.Model) // normally, don't ignore errors!
	if err == nil && camModel != nil {
		fmt.Printf("  Camera Model: %v\n", camModel)
	}
	ownName, _ := x.Get(mknote.OwnerName) // normally, don't ignore errors!
	if err == nil && ownName != nil {
		fmt.Printf("  Owner Name: %v\n", ownName)
	}

	// x.Walk(Walker{})
	return nil
}

type Walker struct{}

func (_ Walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	data, _ := tag.MarshalJSON()
	fmt.Printf("    %v: %v\n", name, string(data))
	return nil
}

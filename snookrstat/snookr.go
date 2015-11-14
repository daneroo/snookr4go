package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/daneroo/snookr4go/fswalker"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"log"
	"os"
	"sort"
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

		err := fswalker.WalkImages(root, visit)
		if err != nil {
			fmt.Printf("WalkImages(%s) returned %v\n", root, err)
		} else {
			fmt.Printf("WalkImages(%s) is done\n", root)
		}

		paths := []string{"coco"}
		sort.Strings(paths)
		// m, err := fswalker.MD5All(os.Args[1])
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// var paths []string
		// for path := range m {
		// 	paths = append(paths, path)
		// }
		// sort.Strings(paths)
		// for _, path := range paths {
		// 	fmt.Printf("%x  %s\n", m[path], path)
		// }

	}

}

func visit(ima fswalker.ImageInfo) error {
	name := ima.FileName
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

	asJson, err := json.Marshal(ima)
	fmt.Printf("  ---- Image '%s' ----\n", asJson)
	stamp, err := x.DateTime() // normally, don't ignore errors!
	if err != nil {
		fmt.Printf("  Date: %v\n", stamp)
	} else {
		fmt.Printf("  Date error: %v\n", err)

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

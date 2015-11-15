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
	"time"
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
		}

		// the bounded walker
		if err := fswalker.MD5All(root); err != nil {
			fmt.Printf("MD5 error:%v\n", err)
		}

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

	taken, err := x.DateTime() // normally, don't ignore errors!
	if err == nil {
		ima.Taken = taken
	} else {
		fmt.Printf("  Date error: %v\n", err)
		ima.Taken = time.Unix(0, 0)
	}

	model, err := x.Get(exif.Model) // normally, don't ignore errors!
	if err == nil && model != nil {
		ima.Camera = fmt.Sprintf("%v", model)
	}

	owner, err := x.Get(mknote.OwnerName) // normally, don't ignore errors!
	if err == nil && owner != nil {
		ima.Owner = fmt.Sprintf("%v", owner)
	}

	asJson, err := json.Marshal(ima)
	fmt.Println(string(asJson))
	// x.Walk(Walker{})

	return nil
}

type Walker struct{}

func (_ Walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	data, _ := tag.MarshalJSON()
	fmt.Printf("    %v: %v\n", name, string(data))
	return nil
}

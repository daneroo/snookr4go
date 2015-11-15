// Package fswalker implements walking a folder and visiting images
// specification (http://www.exif.org/Exif2-2.PDF).
package fswalker

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// for Marshalling without TimeZone
type ZonelessTime time.Time

func (t ZonelessTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05"))
	return []byte(stamp), nil
}

type ImageInfo struct {
	FileName     string       `json:"filename"`
	Size         int64        `json:"size"`
	Md5          string       `json:"md5"`
	LastModified ZonelessTime `json:"lastModified"`
	Taken        ZonelessTime `json:"taken"`
	Camera       string       `json:"camera"`
	Owner        string       `json:"owner"`
}

type ImageWalkFunc func(ima ImageInfo) error

// from filepath: Walk walks the file tree rooted at root.
// calling walkFn for each image (jpg) file
// All errors that arise visiting files and directories are filtered by walkFn.
//  -The files are walked in lexical order,
//  -Walk does not follow symbolic links.
func WalkImages(root string, walkImageFn ImageWalkFunc) error {
	filterFn := func(filename string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error visiting %s: %s\n", filename, err)
		}

		ext := strings.ToLower(path.Ext(filename))
		if !f.IsDir() && (ext == ".jpg" || ext == ".jpeg") {

			// for MD5
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			ima := ImageInfo{
				FileName:     filename,
				Size:         f.Size(),
				LastModified: ZonelessTime(f.ModTime()),
				Md5:          fmt.Sprintf("%x", md5.Sum(data)),
			}

			walkImageFn(ima)
		}

		return nil
	}
	return filepath.Walk(root, filterFn)
}

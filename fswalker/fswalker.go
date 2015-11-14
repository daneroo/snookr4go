// Package fswalker implements walking a folder and visiting images
// specification (http://www.exif.org/Exif2-2.PDF).
package fswalker

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ImageInfo struct {
	FileName     string
	Size         int64
	Md5          string
	LastModified time.Time
	Taken        time.Time
	Camera       string
	Owner        string
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

			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			ima := ImageInfo{
				FileName:     filename,
				Size:         f.Size(),
				LastModified: f.ModTime(),
				Md5:          fmt.Sprintf("%x", md5.Sum(data)),
			}

			walkImageFn(ima)
		}

		return nil
	}
	return filepath.Walk(root, filterFn)
}

// from https://blog.golang.org/pipelines/bounded.go
// walkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel.  If done is closed, walkFiles abandons its work.
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() { // HL
		// Close the paths channel after Walk returns.
		defer close(paths) // HL
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path: // HL
			case <-done: // HL
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// digester reads path names from paths and sends digests of the corresponding
// files on c until either paths or done is closed.
func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
	for path := range paths { // HLpaths
		data, err := ioutil.ReadFile(path)
		select {
		case c <- result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.  In that case,
// MD5All does not wait for inflight read operations to complete.
func MD5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan result) // HLc
	var wg sync.WaitGroup
	const numDigesters = 20
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(done, paths, c) // HLc
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c) // HLc
	}()
	// End of pipeline. OMIT

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		return nil, err
	}
	return m, nil
}
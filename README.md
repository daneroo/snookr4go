# Snookr for Go

Re-implement the parts of snookr4gv2, which is in Groovy

Start from example in [this tutorial](https://github.com/GoesToEleven/GolangTraining/blob/master/50_exif/main.go), 
and also [the example in exifstat folder of rwcarlsen's exif lib](https://github.com/rwcarlsen/goexif/blob/go1/exifstat/main.go)
## The parts

-Scan the file system
-Read Exif data for a file: Date,Camera,Owner
-Get FileInfo, md5, sha1
-[Marshal JSON](http://blog.golang.org/json-and-go)
-Talk to Flickr: getPhotoList, uploadPhoto. (oauth can be reused from snookr4gv2/auth-node)
-Define interfaces for...

## Code Organization

- `snookrstat`: Command line utility
    
    go run snookrstat/snookr.go images/ notExist images/dad.jpg 

- FileWalking
- Exif manipulation
- Flickr stuff

## References
### JSON

- [Golang Blog: JSON](http://blog.golang.org/json-and-go)

### Concurrency

- [Golang Blog: Pipelines](https://blog.golang.org/pipelines)
- [Effective Go](http://golang.org/doc/effective_go.html#concurrency)

### Flickr packages

- [masci](https://github.com/masci/flickr) : seems the way to go!
- [Camlistore](https://camlistore.org/pkg/importer/flickr/flickr.go) : just an importer
- [manki](https://github.com/manki/flickgo) : imcomplete
- [mncaudill](https://github.com/mncaudill/go-flickr) : verryyy old



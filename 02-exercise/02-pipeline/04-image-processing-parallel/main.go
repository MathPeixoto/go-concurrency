package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type fileImage struct {
	path string
	img  *image.NRGBA
}

// Image processing - sequential
// Input - directory with images.
// output - thumbnail images

// in order to parallelize the image processing I will use pipelines
// walkfiles -> processImage -> saveThumbnail
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}
	start := time.Now()

	chPath, err := walkFiles(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// process the image
	chThumbnailImage := processImage(chPath)

	for file := range chThumbnailImage {
		wg.Add(1)
		// save the thumbnail image to disk
		go func(file *fileImage) {
			defer wg.Done()
			err := saveThumbnail(file.path, file.img)
			if err != nil {
				return
			}
		}(file)
	}

	wg.Wait()
	fmt.Printf("Time taken: %s\n", time.Since(start))
}

// walfiles - take diretory path as input
// does the file walk
// generates thumbnail images
// saves the image to thumbnail directory.
func walkFiles(root string) (<-chan string, error) {
	var wg sync.WaitGroup

	chPath := make(chan string)

	fn := func(path string, info os.FileInfo, err error) error {

		// filter out error
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			// check if it is file
			if !info.Mode().IsRegular() {
				return
			}

			fmt.Println("WALK FILES")

			// check if it is image/jpeg
			contentType, _ := getFileContentType(path)
			if contentType == "image/jpeg" {
				chPath <- path
			}

		}(path)

		return nil
	}

	err := filepath.Walk(root, fn)

	if err != nil {
		return nil, err
	}

	go func() {
		wg.Wait()
		close(chPath)
	}()

	return chPath, nil
}

// processImage - takes image file as input
// return pointer to thumbnail image in memory.
func processImage(chPath <-chan string) <-chan *fileImage {
	var wg sync.WaitGroup

	chProcessedImage := make(chan *fileImage)

	for p := range chPath {

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			fmt.Println("PROCESSING IMAGE")

			// load the image from file
			srcImage, err := imaging.Open(path)
			if err != nil {
				return
			}

			// scale the image to 100px * 100px
			thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

			chProcessedImage <- &fileImage{
				path: path,
				img:  thumbnailImage,
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(chProcessedImage)
	}()

	return chProcessedImage
}

// saveThumbnail - save the thumbnail image to folder
func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	//fmt.Println("SAVE THUMBNAIL")

	filename := filepath.Base(srcImagePath)
	//dstImagePath := "thumbnail/" + filename
	dstImagePath := "/home/matheus/GolandProjects/go-concurrency-exercises/02-exercise/02-pipeline/04-image-processing-parallel/thumbnail/" + filename

	// save the image in the thumbnail folder.
	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}

// getFileContentType - return content type and error status
func getFileContentType(file string) (string, error) {

	out, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(out)

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

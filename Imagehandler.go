package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

// please put the file path.
// openfile > decodefile > resizefile > createemptyfile > encodefileto created one > return file.
func Imagehandler(str string) {
	file, err := os.Open(str)
	if err != nil {
		log.Fatal("fatal error occured", err)
	}
	defer file.Close()
	decodedfile, _, _ := image.Decode(file)
	if err != nil {
		log.Fatal("fatal error occured", err)
	}

	resizedimg := resize.Resize(300, 400, decodedfile, resize.Lanczos3)

	resiezedimgfile, _ := os.Create("resizedimage.png")
	defer resiezedimgfile.Close()

	png.Encode(resiezedimgfile, resizedimg)

	fmt.Println("the image that you has input has been encoded with new size png file.")
}

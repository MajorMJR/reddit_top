package main

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	//"log"
	"fmt"
	"os"
)

func resizeImg(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
	}

	m := resize.Resize(800, 0, img, resize.Lanczos3)

	fmt.Println(filename[7:])
	filename = "img/resized/" + filename[4:]
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	return nil
}

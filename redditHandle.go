package main

import (
	"fmt"
	"github.com/jzelinskie/reddit"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	//"reflect"
)

type Image struct {
	URL      string
	Title    string
	Filename string
}

func listImagesDl(subreddit string) []Image {
	session, err := reddit.NewLoginSession("imgdl", "testing123", "GO BOT")
	if err != nil {
		fmt.Println(err)
	}

	submissions, err := session.SubredditSubmissions(subreddit)
	if err != nil {
		fmt.Println(err)
	}

	images := []Image{}

	for _, k := range submissions {
		if k.URL[len(k.URL)-4:] == ".jpg" {
			//fmt.Println("this is a jpeg", k.URL)
			imgFilename := strings.Replace(k.Title, " ", "_", -1)
			image := Image{
				URL:      k.URL,
				Title:    k.Title,
				Filename: imgFilename,
			}
			if err != nil {
				//fmt.Println(err)
			}
			images = append(images, image)
		}

	}

	return images
}

func loadImages() ([]Image, error) {
	images_file := []Image{}

	files, _ := ioutil.ReadDir("./img")
	for _, f := range files {
		imgURL := "img/" + f.Name()
		imgTitle := strings.Replace(f.Name(), "_", " ", -1)
		image := Image{
			URL:      imgURL,
			Title:    imgTitle[3:],
			Filename: f.Name(),
		}
		images_file = append(images_file, image)
	}

	return images_file, nil
}

func downloadImages(w http.ResponseWriter, images []Image, subreddit string) error {

	for _, img := range images {
		imgFilename := "img/dl_" + strings.Replace(img.Title, " ", "_", -1)
		if _, err := os.Stat(imgFilename); err == nil {
			fmt.Println("file already exists: ", imgFilename)
			return err
		}
		output, err := os.Create(imgFilename)

		reqImg, err := http.Get(img.URL)
		if err != nil {
			fmt.Println(err)
		}
		defer reqImg.Body.Close()

		n, err := io.Copy(output, reqImg.Body)
		fmt.Println(n, "bytes downloaded")
	}
	return nil
}

func RedditHandler(w http.ResponseWriter, r *http.Request) {
	//image_files := listImagesDl("earthporn")
	//downloadImages(w, image_files, "earthporn")
	image_files, _ := loadImages()

	//fmt.Println(images)
	t, err := template.ParseFiles("tmpl/reddit.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, map[string]interface{}{
		"Images": image_files,
	})
}

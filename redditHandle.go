package main

import (
	"fmt"
	"github.com/jzelinskie/reddit"
	"html/template"
	"io"
	"net/http"
	"os"
	//"reflect"
)

type Image struct {
	URL   string
	Title string
}

func listImages(subreddit string) []Image {
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
		image := Image{
			URL:   k.URL,
			Title: k.Title,
		}
		images = append(images, image)
	}

	return images
}

func saveImages(w http.ResponseWriter, images []Image, subreddit string) error {
	for _, img := range images {
		filename := "img/dl_" + img.Title
		if _, err := os.Stat(filename); err == nil {
			fmt.Println("file already exists: ", filename)
		}
		output, err := os.Create(filename)

		reqImg, err := http.Get(img.URL)
		if err != nil {
			fmt.Println(err)
		}
		defer reqImg.Body.Close()

		n, err := io.Copy(output, reqImg.Body)
		fmt.Println(n, "bytes downloaded")
	}
}

func RedditHandler(w http.ResponseWriter, r *http.Request) {
	session, err := reddit.NewLoginSession("Injunire", "mitjamroo", "GO BOT")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(reflect.TypeOf(session))
	submissions, _ := session.SubredditSubmissions("earthporn")

	images := listImages("earthporn")

	saveImages(w, images, "earthporn")

	//fmt.Println(images)
	t, err := template.ParseFiles("tmpl/reddit.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, map[string]interface{}{
		"Submissions": submissions,
	})
}

package main

import (
	"fmt"
	"github.com/jzelinskie/reddit"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	//"reflect"
	"strings"
)

type Image struct {
	URL      string
	Title    string
	Filename string
	size800  int
}

type User struct {
	Username  string
	Password  string
	Useragent string
}

func (u *User) loginReddit() (*reddit.LoginSession, error) {
	session, err := reddit.NewLoginSession(u.Username, u.Password, u.Useragent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("user authenticated")
	return session, nil
}

func getSubrdt(subreddit string) {}

func listImagesDl(subreddit string) []Image {
	//session, err := reddit.NewLoginSession("imgdl", "testing123", "GO BOT")
	//if err != nil {
	//	fmt.Println(err)
	//}
	user := User{
		Username:  "imgdl",
		Password:  "testing123",
		Useragent: "GO BOT",
	}
	Session, _ := user.loginReddit()
	submissions, err := Session.SubredditSubmissions(subreddit)
	if err != nil {
		fmt.Println(err)
	}

	images := []Image{}

	for _, k := range submissions {
		if k.URL[len(k.URL)-4:] == ".jpg" {
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
	fmt.Println(images)
	return images
}

func downloadImages(w http.ResponseWriter, images []Image, subreddit string) ([]Image, error) {
	// TODO sort downloaded images into dir organized by subreddit

	for _, img := range images {
		//fmt.Println(img.Filename)
		imgFilename := "img/dl_" + img.Filename
		//if file exists then skip it
		if _, err := os.Stat(imgFilename); err == nil {
			//fmt.Println("file already exists: ")
			continue
		}
		output, err := os.Create(imgFilename)
		fmt.Println(img.URL)
		reqImg, err := http.Get(img.URL)
		if err != nil {
			fmt.Println(err)
		}
		defer reqImg.Body.Close()

		n, err := io.Copy(output, reqImg.Body)
		fmt.Println(n, "bytes downloaded")
	}
	return images, nil
}

func loadImages() ([]Image, error) {
	images_struct := []Image{}

	files, _ := ioutil.ReadDir("./img")
	for _, f := range files {
		imgURL := "img/" + f.Name()
		imgTitle := strings.Replace(f.Name()[3:], "_", " ", -1)
		image := Image{
			URL:      imgURL,
			Title:    imgTitle,
			Filename: f.Name(),
		}
		images_struct = append(images_struct, image)
	}
	resizeImages(images_struct)
	return images_struct, nil
}

func resizeImages(images []Image) error {
	for i, img := range images {
		if _, err := os.Stat("img/resized/" + img.Filename); err == nil {
			fmt.Println("file already exists: ")
			continue
		}
		fmt.Println(len(images), i)
		if i+2 > len(images) {
			fmt.Println("end of images", i)
			return nil
		}
		filename := "img/" + img.Filename

		resizeImg(filename)
	}
	return nil
}

func redditHandler(w http.ResponseWriter, r *http.Request) {

	image_files := listImagesDl("earthporn")
	downloadImages(w, image_files, "earthporn")
	image_files, _ = loadImages()

	//fmt.Println(images)
	t, err := template.ParseFiles("tmpl/reddit.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, map[string]interface{}{
		"Images": image_files,
	})
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	imagefiles, err := loadImages()
	if err != nil {
		fmt.Println(err)
	}
	t, err := template.ParseFiles("tmpl/reddit.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, map[string]interface{}{
		"Images": imagefiles,
	})
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/reddit/", redditHandler)
	http.HandleFunc("/img/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.ListenAndServe(":8080", nil)
}

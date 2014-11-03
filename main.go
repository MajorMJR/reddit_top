package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "pages/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("tmpl/" + tmpl + ".html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, p)
}

func renderList(w http.ResponseWriter, templates []Page) {
	t, err := template.ParseFiles("tmpl/list.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, map[string]interface{}{
		"Templates": templates,
	})
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
		fmt.Println(err)
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir("pages/")
	templates := []Page{}

	for _, f := range files {
		length := len(f.Name())
		filetype := f.Name()[length-3:]
		if filetype == "txt" {
			title := f.Name()[:length-4]
			p, err := loadPage(title)
			if err != nil {
				fmt.Println(err)
			}
			templates = append(templates, *p)
		}
	}
	renderList(w, templates)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/reddit/", redditHandler)
	http.HandleFunc("/img/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	//http.HandleFunc("/img/resized/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, r.URL.Path[1:])
	//})
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Page struct {
	Title      string
	Stylesheet string
	Body       []byte
}

var (
	tmplDir = "tmpl/"
	dataDir = "data/"
)

func (p *Page) save() error {
	filename := dataDir + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title, stylesheet string) (*Page, error) {
	filename := dataDir + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Stylesheet: stylesheet, Body: body}, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "home", &Page{Title: "Homepage", Stylesheet: "home.css"})

}

// func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
// 	p, err := loadPage(title)
// 	if err != nil {
// 		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
// 		return
// 	}
// 	renderTemplate(w, "view", p)
// }

// func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
// 	body := r.FormValue("body")
// 	p := &Page{Title: title, Body: []byte(body)}
// 	err := p.save()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/view/"+title, http.StatusFound)
// }

var validPath = regexp.MustCompile("^/(home|edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Println(m)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {

	fmt.Println("localhost:8080")

	http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/view/", makeHandler(viewHandler))
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.ListenAndServe(":8080", nil)
}

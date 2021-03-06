package main

import (
	"fmt"
	md "github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(edit|save|view|stylesheets)/([a-zA-Z0-9_]+)$")
var validStaticFile = regexp.MustCompile("^/stylesheets/([a-zA-Z0-9_]+).css$")
var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html", "templates/list.html"))

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".md"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".md"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderMarkdown(title string, input []byte) []byte {
	var renderer md.Renderer
	htmlFlags := 0
	extensions := 0
	css := "../stylesheets/styles.css"
	htmlFlags |= md.HTML_USE_XHTML
	htmlFlags |= md.HTML_USE_SMARTYPANTS
	htmlFlags |= md.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= md.HTML_COMPLETE_PAGE
	renderer = md.HtmlRenderer(htmlFlags, title, css)
	return md.Markdown(input, renderer, extensions)

}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func getFiles() ([]string, error) {
	dir, err := os.Open("data")
	if err != nil {
		return nil, err
	}
	fileInfo, err := dir.Readdir(10)
	if err != nil {
		return nil, err
	}
	files := make([]string, len(fileInfo))
	for i, f := range fileInfo {
		files[i] = strings.Split(f.Name(), ".")[0]
	}
	return files, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	p.Body = renderMarkdown(title, p.Body)
	renderTemplate(w, "view", p)
	fmt.Fprintf(w, "%s", p.Body)

}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := getFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = templates.ExecuteTemplate(w, "list.html", files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	m := validStaticFile.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	body, err := ioutil.ReadFile("./" + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", body)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/stylesheets/", staticHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)

}

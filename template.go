package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Header template.HTML
	Navbar template.HTML
	Footer template.HTML
	Custom map[string]string
}

var h template.HTML
var f template.HTML

func ApplyTemplate(FileName string, Inject string, res http.ResponseWriter) {
	b, _ := ioutil.ReadFile(FileName)
	t := template.New("TPage")
	t.Parse(string(b))
	t.Execute(res, Inject)
}

func renderTemplate(fileName string, custom map[string]string, res http.ResponseWriter) {
	nav := getNavbarTemplate(custom["navbarActive"])
	p := &Page{Header: h, Navbar: nav, Footer: f, Custom: custom}
	t, _ := template.ParseFiles(fileName)
	t.Execute(res, p)
}

func getNavbarTemplate(active string) template.HTML {
	custom := map[string]string{
		"search":  "inactive",
		"chart":   "inactive",
		"network": "inactive",
		"map":     "inactive",
	}
	custom[active] = "active"
	p := &Page{Custom: custom}
	t, _ := template.ParseFiles("public/templates/dc-navbar.html")
	var bf bytes.Buffer
	t.Execute(&bf, p)
	return template.HTML(bf.String())
}

func initTemplates() {
	hf, _ := ioutil.ReadFile("public/templates/header.html")
	h = template.HTML(hf)
	ff, _ := ioutil.ReadFile("public/templates/footer.html")
	f = template.HTML(ff)
}

package main

import (
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
var n template.HTML
var f template.HTML

func ApplyTemplate(FileName string, Inject string, res http.ResponseWriter) {
	b, _ := ioutil.ReadFile(FileName)
	t := template.New("TPage")
	t.Parse(string(b))
	t.Execute(res, Inject)
}

func renderTemplate(fileName string, custom map[string]string, res http.ResponseWriter) {
	p := &Page{Header: h, Navbar: n, Footer: f, Custom: custom}
	t, _ := template.ParseFiles(fileName)
	t.Execute(res, p)
}

func initTemplates() {
	hf, _ := ioutil.ReadFile("public/templates/header.html")
	h = template.HTML(hf)
	nf, _ := ioutil.ReadFile("public/templates/navbar.html")
	n = template.HTML(nf)
	ff, _ := ioutil.ReadFile("public/templates/footer.html")
	f = template.HTML(ff)
}

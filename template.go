package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

func ApplyTemplate(FileName string, Inject string, res http.ResponseWriter) {
	b, _ := ioutil.ReadFile(FileName)
	t := template.New("TPage")
	t.Parse(string(b))
	t.Execute(res, Inject)
}

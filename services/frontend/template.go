package main

//
//import (
//	"fmt"
//	"html/template"
//	"io/ioutil"
//	"net/http"
//	"path/filepath"
//)
//
//// parseTemplate applies a given file to the body of the base template.
//func parseTemplate(filename string) *appTemplate {
//	tmpl := template.Must(template.ParseFiles("templates/book_base.html"))
//
//	// Put the named file into a template called "body"
//	path := filepath.Join("templates", filename)
//	b, err := ioutil.ReadFile(path)
//	if err != nil {
//		panic(fmt.Errorf("could not read template: %v", err))
//	}
//	template.Must(tmpl.New("body").Parse(string(b)))
//
//	return &appTemplate{tmpl.Lookup("book_base.html")}
//}
//
//// appTemplate is an appError-aware wrapper for a html/template.
//type appTemplate struct {
//	t *template.Template
//}
//
//// Execute writes the template using the provided data.
//func (tmpl *appTemplate) Execute(fe *frontendServer, w http.ResponseWriter, r *http.Request, data interface{}) *appError {
//	d := struct {
//		Data interface{}
//	}{
//		Data: data,
//	}
//
//	if err := tmpl.t.Execute(w, d); err != nil {
//		return fe.appErrorf(r, err, "could not write template: %v", err)
//	}
//	return nil
//}

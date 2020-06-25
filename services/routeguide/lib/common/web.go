package common

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

// https://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	err     error
	message string
	code    int
	req     *http.Request
	c       *AppConfig
	stack   []byte
}

// parseTemplate applies a given file to the body of the base template.
func parseTemplate(filename string) *appTemplate {
	tmpl := template.Must(template.ParseFiles("templates/base.html"))

	// Put the named file into a template called "body"
	path := filepath.Join("templates", filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		App.Log.Errorf("could not read template: %v", err)
		panic(fmt.Errorf("could not read template: %v", err))
	}
	template.Must(tmpl.New("body").Parse(string(b)))

	return &appTemplate{tmpl.Lookup("base.html")}
}

// appTemplate is an appError-aware wrapper for a html/template.
type appTemplate struct {
	t *template.Template
}

// Execute writes the template using the provided data.
func (tmpl *appTemplate) Execute(c *AppConfig, w http.ResponseWriter, r *http.Request, data interface{}) *appError {
	d := struct {
		Data interface{}
	}{
		Data: data,
	}

	if err := tmpl.t.Execute(w, d); err != nil {
		return c.appErrorf(r, err, "could not write template: %v", err)
	}
	return nil
}

func (c *AppConfig) appErrorf(r *http.Request, err error, format string, v ...interface{}) *appError {
	return &appError{
		err:     err,
		message: fmt.Sprintf(format, v...),
		code:    500,
		req:     r,
		c:       c,
		stack:   debug.Stack(),
	}
}

//
// Based on version from Konstanin Ivanov <kostyarin.ivanov@gmail.com>

// Basically this loads templates and uses the directory structure to provide
// prefix.

// So if you you load all templates in templates and there is an edit.gohtml in subdir then
// the name is subdir_edit

// A Tmpl implements keeper, loader and reloader for HTML templates
type Tmpl struct {
	*template.Template                  // root template
	dir                string           // root directory
	ext                string           // extension
	devel              bool             // reload every time
	funcs              template.FuncMap // functions
	loadedAt           time.Time        // loaded at (last loading time)
}

// NewTmpl creates new Tmpl and loads templates. The dir argument is
// directory to load templates from. The ext argument is extension of
// templates. The devel (if true) turns the Tmpl to reload templates
// every Render if there is a change in the dir.
func NewTmpl(dir, ext string, devel bool) (tmpl *Tmpl, err error) {
	// get absolute path
	if dir, err = filepath.Abs(dir); err != nil {
		return
	}

	tmpl = new(Tmpl)
	tmpl.dir = dir
	tmpl.ext = ext
	tmpl.devel = devel

	if err = tmpl.Load(); err != nil {
		tmpl = nil // drop for GC
	}

	return
}

// Dir returns absolute path to directory with views
func (t *Tmpl) Dir() string {
	return t.dir
}

// Ext returns extension of views
func (t *Tmpl) Ext() string {
	return t.ext
}

// Devel returns development pin
func (t *Tmpl) Devel() bool {
	return t.devel
}

// Funcs sets template functions
func (t *Tmpl) Funcs(funcMap template.FuncMap) {
	t.Template = t.Template.Funcs(funcMap)
	t.funcs = funcMap
}

// Load or reload templates
func (t *Tmpl) Load() (err error) {
	// time point
	t.loadedAt = time.Now()

	// unnamed root template
	var root = template.New("")

	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {
		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != t.ext {
			return
		}

		// get relative path
		var rel string
		if rel, err = filepath.Rel(t.dir, path); err != nil {
			return err
		}

		// name of a template is its relative path
		// without extension
		rel = strings.TrimSuffix(rel, t.ext)

		// load or reload
		var (
			nt = root.New(rel)
			b  []byte
		)

		if b, err = ioutil.ReadFile(path); err != nil {
			return err
		}

		_, err = nt.Parse(string(b))
		return err
	}

	if err = filepath.Walk(t.dir, walkFunc); err != nil {
		return
	}

	// necessary for reloading
	if t.funcs != nil {
		root = root.Funcs(t.funcs)
	}

	t.Template = root // set or replace
	return
}

// IsModified lookups directory for changes to
// reload (or not to reload) templates if development
// pin is true.
func (t *Tmpl) IsModified() (yep bool, err error) {

	var errStop = errors.New("stop")

	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {

		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != t.ext {
			return
		}

		if yep = info.ModTime().After(t.loadedAt); yep == true {
			return errStop
		}

		return
	}

	// clear the errStop
	if err = filepath.Walk(t.dir, walkFunc); err == errStop {
		err = nil
	}

	return
}

func (t *Tmpl) Render(w io.Writer, name string, data interface{}) (err error) {

	// if devlopment
	if t.devel == true {

		// lookup directory for changes
		var modified bool
		if modified, err = t.IsModified(); err != nil {
			return
		}

		// reload
		if modified == true {
			if err = t.Load(); err != nil {
				return
			}
		}

	}

	err = t.ExecuteTemplate(w, name, data)
	return
}

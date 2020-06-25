package main

import (
	"cloud.google.com/go/storage"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"io/ioutil"
	"lib/common"
	"net/http"
	"os"
	"path"

	pb "frontend/pb/pb_book_v1"
)

var (
	bookTemplates = template.Must(template.New("").ParseGlob("templates/book/*.gohtml"))
	ErrNeedBookID = errors.New("Need a book ID")
)

// listHandler displays a list of books in the database.
func (fe *frontendServer) listBook(w http.ResponseWriter, r *http.Request) {
	fe.log.Debug("List books")
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	ctx := r.Context()
	books, err := fe.ListBooks(ctx)
	if err != nil {
		renderHTTPError(r, w, fmt.Errorf("could not retrieve books. %w", err))
		return
	}
	if err := bookTemplates.ExecuteTemplate(w, "list", map[string]interface{}{
		"session_id":    sessionID(r),
		"request_id":    r.Context().Value(ctxKeyRequestID{}),
		"banner_color":  common.App.CanaryColour, // illustrates canary deployments
		"books":         books,
		"platform_url":  common.App.Platform.Url,
		"platform_name": common.App.Platform.Provider,
	}); err != nil {
		log.Println(err)
	}
}

// addBook displays a blank edit form that captures details of a new book to add
func (fe *frontendServer) addBook(w http.ResponseWriter, r *http.Request) {
	fe.log.Debug("Add Book")
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	if err := bookTemplates.ExecuteTemplate(w, "edit", map[string]interface{}{
		"session_id":    sessionID(r),
		"request_id":    r.Context().Value(ctxKeyRequestID{}),
		"banner_color":  common.App.CanaryColour, // illustrates canary deployments
		"platform_url":  common.App.Platform.Url,
		"platform_name": common.App.Platform.Provider,
	}); err != nil {
		log.Println(err)
	}
}

// bookDetail displays the details of a given book.
func (fe *frontendServer) bookDetail(w http.ResponseWriter, r *http.Request) {
	fe.log.Debug("Book Details")
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	book, err := fe.bookFromRequest(r)
	if err != nil {
		renderHTTPError(r, w, fmt.Errorf("could not retrieve books. %w", err))
		return
	}
	if err := bookTemplates.ExecuteTemplate(w, "detail.gohtml", map[string]interface{}{
		"session_id":    sessionID(r),
		"request_id":    r.Context().Value(ctxKeyRequestID{}),
		"banner_color":  common.App.CanaryColour, // illustrates canary deployments
		"book":          book,
		"platform_url":  common.App.Platform.Url,
		"platform_name": common.App.Platform.Provider,
	}); err != nil {
		log.Println(err)
	}
}

// bookFromRequest retrieves a book given a book ID in the URL's path.
func (fe *frontendServer) bookFromRequest(r *http.Request) (*pb.Book, error) {
	fe.log.Debug("Book from request")
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, ErrNeedBookID
	}
	book, err := fe.GetBook(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find book: %w", err)
	}
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Info("Read from service", book)
	return book, nil
}

// editBook shows the details of a given book to edit.
func (fe *frontendServer) editBook(w http.ResponseWriter, r *http.Request) {
	fe.log.Debug("Edit book")
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	book, err := fe.bookFromRequest(r)
	if err != nil {
		renderHTTPError(r, w, fmt.Errorf("could not retrieve book. %w", err))
		return
	}
	if err := bookTemplates.ExecuteTemplate(w, "edit", map[string]interface{}{
		"session_id":    sessionID(r),
		"request_id":    r.Context().Value(ctxKeyRequestID{}),
		"banner_color":  common.App.CanaryColour, // illustrates canary deployments
		"platform_url":  common.App.Platform.Url,
		"platform_name": common.App.Platform.Provider,
		"book":          book,
	}); err != nil {
		log.Println(err)
	}
}

// bookFromForm populates the fields of a Book from form values
// (see templates/book/edit.gohtml).
func (fe *frontendServer) bookFromForm(r *http.Request) (*pb.Book, error) {
	fe.log.Debug("Book From Form")
	//fe.log.Debug(r.ParseForm())
	//fe.log.Debug(r.Form)
	//ctx := r.Context()
	imageURL := ""
	imageURL, err := fe.uploadFileFromForm(r)
	if err != nil {
		return nil, fmt.Errorf("could not upload file: %w", err)
	}
	if imageURL == "" {
		imageURL = r.FormValue("imageURL")
	}
	// Get the book details
	book := &pb.Book{
		Title:         r.FormValue("title"),
		Author:        r.FormValue("author"),
		PublishedDate: r.FormValue("publishedDate"),
		ImageURL:      imageURL,
		Description:   r.FormValue("description"),
	}
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Info("Read from form:", book)
	return book, nil
}

// createBook adds a book
func (fe *frontendServer) createBook(w http.ResponseWriter, r *http.Request) {
	fe.log.Debug("Create Book")
	ctx := r.Context()
	book, err := fe.bookFromForm(r)
	if err != nil {
		renderHTTPError(r, w, fmt.Errorf("could not parse book from form: %w", err))
	}
	id, err := fe.AddBook(ctx, book)
	if err != nil {
		renderHTTPError(r, w, fmt.Errorf("could not save book from form: %w", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/books/%s", id), http.StatusFound)
	//w.Header().Set("location", "/cart")
	//w.WriteHeader(http.StatusFound)
}

// updateBook updates the details of a given book.
func (fe *frontendServer) updateBook(w http.ResponseWriter, r *http.Request) {
	fe.log.Debug("Update book")
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	if id == "" {
		fe.log.Errorf("Cannot update book: %v", ErrNeedBookID)
		renderHTTPError(r, w, fmt.Errorf("Cannot update book: %w", ErrNeedBookID))
	}
	book, err := fe.bookFromForm(r)
	if err != nil {
		fe.log.Errorf("could not update book from form: %v", err)
		renderHTTPError(r, w, fmt.Errorf("could not update book from form: %w", err))
	}
	book.Id = id

	if book, err = fe.UpdateBook(ctx, book); err != nil {
		renderHTTPError(r, w, fmt.Errorf("could not update book: %w", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/books/%s", book.Id), http.StatusFound)
}

// deleteBook deletes a given book.
func (fe *frontendServer) deleteBook(w http.ResponseWriter, r *http.Request) {
	fe.log.Debug("Delete book")
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	if err := fe.DeleteBook(ctx, id); err != nil {
		renderHTTPError(r, w, fmt.Errorf("could not update book: %w", err))
	}
	http.Redirect(w, r, "/books", http.StatusFound)
}

func (fe frontendServer) bookTemplates(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", bookTemplates)
}

// uploadFileFromForm uploads a file if it's present in the "image" form field.
//This routing uploads the file to a storage bucket and then provides a URL to the file
//very nice - need to think about how to make this generic and platform specific
func (fe *frontendServer) uploadFileFromForm(r *http.Request) (url string, err error) {
	ctx := r.Context()
	if r.FormValue("image") == "" {
		return "", nil
	}
	fe.log.Debugf("image Form Value %v", r.FormValue("image"))
	f, fh, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		fe.log.Error("Error missing file:", err)
		return "", nil
	}
	if err != nil {
		fe.log.Errorf("Error retrieving the file:%v", err)
		return "", err
	}
	defer f.Close()
	fe.log.Infof("Uploaded File: %+v", fh.Filename)
	fe.log.Infof("File Size: %+v", fh.Size)
	fe.log.Infof("MIME Header: %+v", fh.Header)
	// random filename, retaining existing extension.
	name := uuid.Must(uuid.NewUUID()).String() + path.Ext(fh.Filename)

	if fe.StorageBucket == nil {
		// Create a temporary file within the temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("temp-images", name)
		if err != nil {
			fe.log.Error("Error uploading file:", err)
			return "", err
		}
		defer tempFile.Close()
		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			fe.log.Error("Error readinf file:", err)
			return "", err
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)
		var imageURL = fmt.Sprintf(os.TempDir()+"/%s/%s", "temp-images", tempFile)
		fe.log.Infof("Image URL: %s", imageURL)
		return imageURL, nil

		//return "", errors.New("storage bucket is missing: check bookshelf.go")
	}
	if _, err := fe.StorageBucket.Attrs(ctx); err != nil {
		if err == storage.ErrBucketNotExist {
			return "", fmt.Errorf("bucket %q does not exist: check bookshelf.go", fe.StorageBucketName)
		}
		return "", fmt.Errorf("could not get bucket: %v", err)
	}

	w := fe.StorageBucket.Object(name).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, fe.StorageBucketName, name), nil
}

// https://play.golang.org/p/ty9SaB7p0eq
//package main
//
//import (
//"html/template"
//"os"
//)
//
//type Book struct {
//	ID            string
//	Title         string
//	Author        string
//	PublishedDate string
//	ImageURL      string
//	Description   string
//}
//
//
//func main() {
//	t := template.Must(template.New("").Parse(src))
//	b := Book{ID:"1",Title:"The mystery of the map and struct",Author:"Gopher",PublishedDate:"2015",ImageURL:"",Description:"descr"}
//
//	params := map[string]interface{}{
//		"session_id":       "sessionID",
//		"request_id":       "r.Context()",
//		"book":		     b,
//		"platform_css":     "CSS",
//		"platform_name":    "Platform.Provider",
//	}
//
//	if err := t.Execute(os.Stdout, params); err != nil {
//		panic(err)
//	}
//}
//
//const src = `{{ . }}
//
//    {{with index $ "book"}}
//      {{.Title}} by {{.Author}} written {{.PublishedDate}}
//    {{end}}`

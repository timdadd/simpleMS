package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"lib/common"
	"net/http"
	"time"
)

var (
	templates = template.Must(template.New("").
		Funcs(template.FuncMap{
			"renderTime": renderTime,
		}).ParseGlob("templates/*.gohtml"))
)

func (fe *frontendServer) logoutHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Debug("logging out")
	for _, c := range r.Cookies() {
		c.Expires = time.Now().Add(-time.Hour * 24 * 365)
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

func renderHTTPError(r *http.Request, w http.ResponseWriter, err error) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Errorf("Rendering the error")
	statusCode := http.StatusInternalServerError
	log.WithField("error", err).Error("request error")
	errMsg := fmt.Sprintf("%+v", err)

	w.WriteHeader(statusCode)
	templates.ExecuteTemplate(w, "error", map[string]interface{}{
		"session_id":    sessionID(r),
		"request_id":    r.Context().Value(ctxKeyRequestID{}),
		"error":         errMsg,
		"status_code":   statusCode,
		"status":        http.StatusText(statusCode),
		"banner_color":  common.App.CanaryColour, // illustrates canary deployments
		"platform_url":  common.App.Platform.Url,
		"platform_name": common.App.Platform.Provider,
	})
}

func sessionID(r *http.Request) string {
	v := r.Context().Value(ctxKeySessionID{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func renderTime() string {
	return time.Now().Format(time.RFC1123)
}

func (fe frontendServer) version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", version)
}

//// logWriter is used for request logging and can be overridden for tests.
////
//// See https://cloud.google.com/logging/docs/setup/go for how to use the
//// Stackdriver logging client. Output to stdout and stderr is automaticaly
//// sent to Stackdriver when running on App Engine.
//logWriter io.Writer

//// sendLog logs a message.
////
//// See https://cloud.google.com/logging/docs/setup/go for how to use the
//// Stackdriver logging client. Output to stdout and stderr is automaticaly
//// sent to Stackdriver when running on App Engine.
//func (b *Bookshelf) sendLog(w http.ResponseWriter, r *http.Request) *appError {
//	fmt.Fprintln(b.logWriter, "Hey, you triggered a custom log entry. Good job!")
//
//	fmt.Fprintln(w, `<html>Log sent! Check the <a href="http://console.cloud.google.com/logs">logging section of the Cloud Console</a>.</html>`)
//
//	return nil
//}

//// sendError triggers an error that is sent to Error Reporting.
//func (b *Bookshelf) sendError(w http.ResponseWriter, r *http.Request) *appError {
//	msg := `<html>Logging an error. Check <a href="http://console.cloud.google.com/errors">Error Reporting</a> (it may take a minute or two for the error to appear).</html>`
//	err := errors.New("uh oh! an error occurred")
//	return b.appErrorf(r, err, msg)
//}
//
//// https://blog.golang.org/error-handling-and-go
//type appHandler func(http.ResponseWriter, *http.Request) *appError
//
//type appError struct {
//	err     error
//	message string
//	code    int
//	req     *http.Request
//	b       *Bookshelf
//	stack   []byte
//}

//func (fe *frontendServer)  systemserver(w http.ResponseWriter, r *http.Request) {
//	i := &common.ServerInstance{}
//	resp, err := client.Do(req)
//	if err != nil {
//		w.WriteHeader(http.StatusServiceUnavailable)
//		fmt.Fprintf(w, "Error: %s\n", err.Error())
//		return
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(w, "Error: %s\n", err.Error())
//		return
//	}
//	err = json.Unmarshal([]byte(body), i)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(w, "Error: %s\n", err.Error())
//		return
//	}
//	tpl.Execute(w, i)
//})

//func (fe *frontendServer) healthz(w http.ResponseWriter, r *http.Request) {
//	resp, err := client.Do(req)
//	if err != nil {
//		w.WriteHeader(http.StatusServiceUnavailable)
//		fmt.Fprintf(w, "Backend could not be connected to: %s", err.Error())
//		return
//	}
//	defer resp.Body.Close()
//	ioutil.ReadAll(resp.Body)
//	w.WriteHeader(http.StatusOK)
//})

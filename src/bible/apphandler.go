package bible

import (
    "net/http"
    "reflect"
    "runtime/debug"

    "appengine"
)

type appError struct {
    code int
}

type appHandler func(http.ResponseWriter, *http.Request)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if e := recover(); e != nil {
            code := http.StatusInternalServerError
            if reflect.TypeOf(e).String() == "*bible.appError" {
                code = e.(*appError).code
            }
            serveError(w, r, code)
        }
    }()

    fn(w, r)
}

func serveError(w http.ResponseWriter, r *http.Request, code int) {
    c := appengine.NewContext(r)
    if code == http.StatusInternalServerError {
        c.Errorf("Stack trace:\n%s", debug.Stack())
    }
    w.WriteHeader(code)
}

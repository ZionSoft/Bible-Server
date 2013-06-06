package bible

import (
    "net/http"

    "appengine"
)

type appError struct {
    code    int
    message string
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if e := recover(); e != nil {
            serveError(w, r, http.StatusInternalServerError, e)
        }
    }()

    if err := fn(w, r); err != nil {
        serveError(w, r, err.code, err.message)
    }
}

func serveError(w http.ResponseWriter, r *http.Request, code int, message interface{}) {
    c := appengine.NewContext(r)
    c.Errorf("Error code: %d, message: %v", code, message)
    w.WriteHeader(code)
}

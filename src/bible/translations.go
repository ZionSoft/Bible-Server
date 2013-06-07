package bible

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strconv"

    "appengine"
    "appengine/blobstore"
    "appengine/datastore"
)

func downloadTranslationHandler(w http.ResponseWriter, r *http.Request) *appError {
    if r.Method != "GET" {
        return &appError{http.StatusMethodNotAllowed, fmt.Sprintf("downloadTranslationHandler: Method '%s' not allowed.", r.Method)}
    }

    blobKey := r.FormValue("blobKey")
    if len(blobKey) == 0 {
        return &appError{http.StatusBadRequest, string("downloadTranslationHandler: Missing parameter 'blobKey'.")}
    }

    blobstore.Send(w, appengine.BlobKey(blobKey))
    return nil
}

func queryTranslationsHandler(w http.ResponseWriter, r *http.Request) *appError {
    if r.Method != "GET" {
        return &appError{http.StatusMethodNotAllowed, fmt.Sprintf("queryTranslationHandler: Method '%s' not allowed.", r.Method)}
    }

    // parses query parameters
    params, err := url.ParseQuery(r.URL.RawQuery)
    if err != nil {
        return &appError{http.StatusBadRequest, fmt.Sprintf("queryTranslationHandler: Malformed query string '%s'.", r.URL.RawQuery)}
    }

    since, err := strconv.ParseInt(params.Get("since"), 10, 32)
    if err != nil || since < 0 {
        since = 0
    }

    offset, err := strconv.ParseInt(params.Get("offset"), 10, 32)
    if err != nil || offset < 0 {
        offset = 0
    }

    limit, err := strconv.ParseInt(params.Get("limit"), 10, 32)
    if err != nil || limit <= 0 || limit > 100 {
        limit = 100
    }

    language := params.Get("language")

    // makes the query
    c := appengine.NewContext(r)
    q := datastore.NewQuery("TranslationInfo").Offset(int(offset)).Limit(int(limit))
    if since > 0 {
        q = q.Filter("Timestamp >=", since)
    }
    if len(language) > 0 {
        q = q.Filter("Language =", language)
    }
    i := q.Run(c)
    translations := make([]*TranslationInfo, 0, limit)
    for {
        translationInfo := new(TranslationInfo)
        key, err := i.Next(translationInfo)
        if err == datastore.Done {
            break
        }
        if err != nil {
            return &appError{http.StatusInternalServerError, fmt.Sprintf("queryTranslationHandler: Failed to read translation info from datastore '%s'.", err.Error())}
        }
        translationInfo.UniqueId = key.IntID()
        translations = append(translations, translationInfo)
    }

    // writes the response
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    buf, _ := json.Marshal(translations)
    fmt.Fprint(w, string(buf))
    return nil
}

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

func downloadTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    blobKey := r.FormValue("blobKey")
    if len(blobKey) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    blobstore.Send(w, appengine.BlobKey(blobKey))
}

func queryTranslationsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // parses query parameters
    params, err := url.ParseQuery(r.URL.RawQuery)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
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
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        translationInfo.UniqueId = key.IntID()
        translations = append(translations, translationInfo)
    }

    // writes the response
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    buf, _ := json.Marshal(translations)
    fmt.Fprint(w, string(buf))
}

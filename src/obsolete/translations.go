/*
 * Copyright (c) 2014 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package obsolete

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strconv"

    "appengine"
    "appengine/blobstore"
    "appengine/datastore"
    "appengine/delay"
    "appengine/memcache"

    "src/core"
)

func DownloadTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    blobKey := r.FormValue("blobKey")
    if len(blobKey) == 0 {
        panic(&core.Error{http.StatusBadRequest, ""})
    }

    // TODO checks if the blob key exists

    // updates logs
    var logTranslationDownloadFunc = delay.Func("logTranslationDownload", logTranslationDownload)
    logTranslationDownloadFunc.Call(appengine.NewContext(r), blobKey)

    // sends the blob
    blobstore.Send(w, appengine.BlobKey(blobKey))
}

func QueryTranslationsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    // fetches all the translation info
    c := appengine.NewContext(r)
    var translations []*TranslationInfo
    memcache.Gob.Get(c, "TranslationInfo", &translations)
    if len(translations) == 0 {
        // missed in memcache, fetches from datastore
        q := datastore.NewQuery("TranslationInfo").Order("Language")
        keys, err := q.GetAll(c, &translations)
        if err != nil {
            panic(&core.Error{http.StatusInternalServerError, err.Error()})
        }
        for i, t := range translations {
            t.UniqueId = keys[i].IntID()
        }

        // updates memcache
        item := &memcache.Item{
            Key:    "TranslationInfo",
            Object: translations,
        }
        memcache.Gob.Set(c, item)
    }

    // parses query parameters
    params, err := url.ParseQuery(r.URL.RawQuery)
    if err != nil {
        panic(&core.Error{http.StatusBadRequest, ""})
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

    // filters translation info based on the query parameters
    filteredTranslations := make([]*TranslationInfo, 0, limit)
    for _, t := range translations {
        if t.Timestamp < since || (len(language) > 0 && t.Language != language) {
            continue
        }

        if offset > 0 {
            offset = offset - 1
            continue
        }

        if limit > 0 {
            limit = limit - 1
        } else {
            break
        }

        filteredTranslations = append(filteredTranslations, t)
    }

    // writes the response
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    buf, _ := json.Marshal(filteredTranslations)
    fmt.Fprint(w, string(buf))
}

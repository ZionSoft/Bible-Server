/*
 * Copyright (c) 2014 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package bible

import (
    "archive/zip"
    "encoding/json"
    "net/http"
    "time"

    "appengine"
    "appengine/blobstore"
    "appengine/datastore"
    "appengine/memcache"

    "src/core"
)

func UploadTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    // reads the uploaded archive
    f, _, err := r.FormFile("translation")
    if err != nil {
        panic(&core.Error{http.StatusBadRequest, ""})
    }
    defer f.Close()

    buf := make([]byte, r.ContentLength)
    n, err := f.Read(buf)
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }
    buf = buf[:n]

    // gets part of the translation info from books.json
    var translationInfo TranslationInfo
    translationInfo.Timestamp = time.Now().Unix()

    reader, err := zip.NewReader(f, int64(n))
    if err != nil {
        panic(&core.Error{http.StatusBadRequest, ""})
    }
    translationInfo.Size = int64(n)

    for _, file := range reader.File {
        if file.Name != "books.json" {
            continue
        }

        rc, err := file.Open()
        if err != nil {
            panic(&core.Error{http.StatusInternalServerError, err.Error()})
        }
        defer rc.Close()
        b := make([]byte, 4096) // 4096 bytes should be big enough
        n, err = rc.Read(b)
        if err != nil {
            panic(&core.Error{http.StatusInternalServerError, err.Error()})
        }
        b = b[:n]

        err = json.Unmarshal(b, &translationInfo)
        if err != nil {
            panic(&core.Error{http.StatusBadRequest, ""})
        }
        break
    }

    // saves uploaded translation into blobstore
    c := appengine.NewContext(r)
    writer, err := blobstore.Create(c, "application/zip")
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, ""})
    }

    _, err = writer.Write(buf)
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }
    err = writer.Close()
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }

    // gets rest of the translation info and writes to datastore
    translationInfo.BlobKey, err = writer.Key()
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }

    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "TranslationInfo", nil), &translationInfo)
    if err != nil {
        blobstore.Delete(c, translationInfo.BlobKey)
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }

    // flushes memcache
    memcache.Flush(c)
}

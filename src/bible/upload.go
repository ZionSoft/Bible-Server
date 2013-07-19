/*
 * Copyright (c) 2013 ZionSoft. All rights reserved.
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
)

func uploadTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        panic(&appError{http.StatusMethodNotAllowed})
    }

    // reads the uploaded archive
    f, _, err := r.FormFile("translation")
    if err != nil {
        panic(&appError{http.StatusBadRequest})
    }
    defer f.Close()

    buf := make([]byte, r.ContentLength)
    n, err := f.Read(buf)
    if err != nil {
        panic(&appError{http.StatusInternalServerError})
    }
    buf = buf[:n]

    // gets part of the translation info from books.json
    var translationInfo TranslationInfo
    translationInfo.Timestamp = time.Now().Unix()

    reader, err := zip.NewReader(f, int64(n))
    if err != nil {
        panic(&appError{http.StatusBadRequest})
    }
    translationInfo.Size = int64(n)

    for _, file := range reader.File {
        if file.Name != "books.json" {
            continue
        }

        rc, err := file.Open()
        if err != nil {
            panic(&appError{http.StatusInternalServerError})
        }
        defer rc.Close()
        b := make([]byte, 4096) // 4096 bytes should be big enough
        n, err = rc.Read(b)
        if err != nil {
            panic(&appError{http.StatusInternalServerError})
        }
        b = b[:n]

        err = json.Unmarshal(b, &translationInfo)
        if err != nil {
            panic(&appError{http.StatusBadRequest})
        }
        break
    }

    // saves uploaded translation into blobstore
    c := appengine.NewContext(r)
    writer, err := blobstore.Create(c, "application/zip")
    if err != nil {
        panic(&appError{http.StatusInternalServerError})
    }

    _, err = writer.Write(buf)
    if err != nil {
        panic(&appError{http.StatusInternalServerError})
    }
    err = writer.Close()
    if err != nil {
        panic(&appError{http.StatusInternalServerError})
    }

    // gets rest of the translation info and writes to datastore
    translationInfo.BlobKey, err = writer.Key()
    if err != nil {
        panic(&appError{http.StatusInternalServerError})
    }

    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "TranslationInfo", nil), &translationInfo)
    if err != nil {
        blobstore.Delete(c, translationInfo.BlobKey)
        panic(&appError{http.StatusInternalServerError})
    }

    // flushes memcache
    memcache.Flush(c)
}

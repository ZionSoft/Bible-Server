/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package translation

import (
    "archive/zip"
    "encoding/json"
    "html/template"
    "net/http"

    "appengine"
    "appengine/blobstore"
    "appengine/datastore"
    "appengine/memcache"

    "src/core"
)

func UploadTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    c := appengine.NewContext(r)
    uploadURL, err := blobstore.UploadURL(c, "/admin/translation/onUploaded", nil)
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }

    w.Header().Set("Content-Type", "text/html")
    err = uploadTemplate.Execute(w, uploadURL)
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }
}

var uploadTemplate = template.Must(template.New("upload").Parse(uploadTemplateHTML))

const uploadTemplateHTML = `
<!--
Copyright (c) 2015 ZionSoft. All rights reserved.
Use of this source code is governed by a BSD-style license
that can be found in the LICENSE file.
-->
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
<html>
<head>
  <title>ZionSoft</title>
  <link rel="icon" type="image/x-icon" href="/view/favicon.ico" />
</head>
<body>
<h2>Upload a New Translation</h2>
<form action="{{.}}" method="post" enctype="multipart/form-data">
  <div><input type="file" name="translation" size="40" /></div>
  <div><input type="submit" name="submit" value="Submit" /></div>
</form>
</body>
</html>
`

func OnTranslationUploadedHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    blobs, _, err := blobstore.ParseUpload(r)
    if err != nil {
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }

    translationBlobs := blobs["translation"]
    if len(translationBlobs) == 0 {
        panic(&core.Error{http.StatusBadRequest, ""})
    }
    translationBlob := translationBlobs[0]

    var translationInfo TranslationInfo
    translationInfo.BlobKey = translationBlob.BlobKey
    translationInfo.Size = translationBlob.Size
    translationInfo.Created = translationBlob.CreationTime.Unix()
    translationInfo.Modified = translationInfo.Created

    c := appengine.NewContext(r)
    blobReader := blobstore.NewReader(c, translationBlob.BlobKey)
    reader, err := zip.NewReader(blobReader, translationBlob.Size)
    if err != nil {
        blobstore.Delete(c, translationBlob.BlobKey)
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }
    for _, file := range reader.File {
        if file.Name == "books.json" {
            rc, err := file.Open()
            if err != nil {
                blobstore.Delete(c, translationBlob.BlobKey)
                panic(&core.Error{http.StatusInternalServerError, err.Error()})
            }
            defer rc.Close()

            b := make([]byte, 4096) // 4096 bytes should be big enough
            n, err := rc.Read(b)
            if err != nil {
                blobstore.Delete(c, translationBlob.BlobKey)
                panic(&core.Error{http.StatusInternalServerError, err.Error()})
            }
            b = b[:n]

            err = json.Unmarshal(b, &translationInfo)
            if err != nil {
                blobstore.Delete(c, translationBlob.BlobKey)
                panic(&core.Error{http.StatusBadRequest, ""})
            }
        } else {
            // TODO validates other files and prepare for full text search
        }
    }

    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "TranslationInfo", nil), &translationInfo)
    if err != nil {
        blobstore.Delete(c, translationBlob.BlobKey)
        panic(&core.Error{http.StatusInternalServerError, err.Error()})
    }

    // flushes memcache
    translations = translations[:0]
    memcache.Flush(c)

    // TODO redirects
}

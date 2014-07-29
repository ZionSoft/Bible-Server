/*
 * Copyright (c) 2014 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package translation

import (
    "net/http"
    "net/url"

    "appengine"
    "appengine/blobstore"

    "src/core"
)

func DownloadTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    // parses query parameters
    params, err := url.ParseQuery(r.URL.RawQuery)
    if err != nil {
        panic(&core.Error{http.StatusBadRequest, ""})
    }

    // TODO supports other query params

    blobKey := params.Get("blobKey")

    translations := loadTranslations(appengine.NewContext(r))
    for _, t := range translations {
        if (string)(t.BlobKey) == blobKey {
            blobstore.Send(w, appengine.BlobKey(blobKey))
            return
        }
    }

    panic(&core.Error{http.StatusBadRequest, ""})
}
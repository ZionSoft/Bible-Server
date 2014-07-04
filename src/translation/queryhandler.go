/*
 * Copyright (c) 2014 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package translation

import (
    "encoding/json"
    "fmt"
    "net/http"

    "appengine"
    "appengine/datastore"
    "appengine/memcache"

    "src/core"
)

func QueryTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    // loads all translations into memory
    c := appengine.NewContext(r)
    var translations []*TranslationInfo
    memcache.Gob.Get(c, "TranslationInfoV2", &translations)
    if len(translations) == 0 {
        // missed memcache, loads from datastore
        q := datastore.NewQuery("TranslationInfoV2")
        keys, err := q.GetAll(c, &translations)
        if err != nil {
            panic(&core.Error{http.StatusInternalServerError, err.Error()})
        }
        for i, t := range translations {
            t.UniqueId = keys[i].IntID()
        }

        // updates memcache
        item := &memcache.Item{
            Key:    "TranslationInfoV2",
            Object: translations,
        }
        memcache.Gob.Set(c, item)
    }

    // TODO supports queries

    // writes the response
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    buf, _ := json.Marshal(translations)
    fmt.Fprint(w, string(buf))
}

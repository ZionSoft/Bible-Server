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

    "src/core"
)

func QueryTranslationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        panic(&core.Error{http.StatusMethodNotAllowed, ""})
    }

    // loads all translations into memory
    c := appengine.NewContext(r)
    translations := loadTranslations(c)

    // TODO supports queries

    // writes the response
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    if translations == nil || len(translations) == 0 {
        fmt.Fprint(w, "[]")
        return
    }

    buf, _ := json.Marshal(translations)
    fmt.Fprint(w, string(buf))
}

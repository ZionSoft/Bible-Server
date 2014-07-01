/*
 * Copyright (c) 2014 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package src

import (
    "net/http"

    "src/bible"
    "src/core"
)

func init() {
    http.Handle("/1.0/downloadTranslation", core.Handler(bible.DownloadTranslationHandler))
    http.Handle("/1.0/translations", core.Handler(bible.QueryTranslationsHandler))

    http.Handle("/admin/uploadTranslation", core.Handler(bible.UploadTranslationHandler))

    http.Handle("/", core.Handler(defaultHandler))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "http://www.zionsoft.net", http.StatusNotFound)
}

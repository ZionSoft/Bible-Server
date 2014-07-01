/*
 * Copyright (c) 2014 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package src

import (
    "net/http"

    "src/core"
    "src/obsolete"
)

func init() {
    http.Handle("/1.0/downloadTranslation", core.Handler(obsolete.DownloadTranslationHandler))
    http.Handle("/1.0/translations", core.Handler(obsolete.QueryTranslationsHandler))

    http.Handle("/admin/uploadTranslation", core.Handler(obsolete.UploadTranslationHandler))

    http.Handle("/", core.Handler(defaultHandler))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "http://www.zionsoft.net", http.StatusNotFound)
}

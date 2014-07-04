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
    "src/translation"
)

func init() {
    http.Handle("/v2/translations", core.Handler(translation.QueryTranslationHandler))
    http.Handle("/admin/translation", core.Handler(translation.UploadTranslationHandler))
    http.Handle("/admin/translation/onUploaded", core.Handler(translation.OnTranslationUploadedHandler))

    // obsoleted
    http.Handle("/1.0/downloadTranslation", core.Handler(obsolete.DownloadTranslationHandler))
    http.Handle("/1.0/translations", core.Handler(obsolete.QueryTranslationsHandler))

    // default handler
    http.Handle("/", core.Handler(defaultHandler))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "http://www.zionsoft.net", http.StatusNotFound)
}

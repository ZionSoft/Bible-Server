/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package src

import (
    "net/http"

    "src/core"
    "src/translation"
)

func init() {
    http.Handle("/v1/translations", core.Handler(translation.QueryTranslationHandler))
    http.Handle("/v1/translation", core.Handler(translation.DownloadTranslationHandler))

    http.Handle("/admin/translation", core.Handler(translation.UploadTranslationHandler))
    http.Handle("/admin/translation/onUploaded", core.Handler(translation.OnTranslationUploadedHandler))

    http.Handle("/", core.Handler(defaultHandler))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "http://www.zionsoft.net", http.StatusNotFound)
}

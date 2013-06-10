/*
 * Copyright (c) 2013 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package bible

import (
    "net/http"
)

func init() {
    http.Handle("/1.0/downloadTranslation", appHandler(downloadTranslationHandler))
    http.Handle("/1.0/translations", appHandler(queryTranslationsHandler))

    http.Handle("/admin/uploadTranslation", appHandler(uploadTranslationHandler))
}

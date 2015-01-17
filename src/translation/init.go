/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package translation

import (
	"net/http"

	"src/core"
)

func init() {
	http.Handle("/v1/translations", core.Handler(queryTranslationHandler))
	http.Handle("/v1/translation", core.Handler(downloadTranslationHandler))

	http.Handle("/admin/translation", core.Handler(uploadTranslationHandler))
	http.Handle("/admin/translation/onUploaded", core.Handler(onTranslationUploadedHandler))
}

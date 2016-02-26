/*
 * Copyright (c) 2016 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package core

import (
	"net/http"
)

func init() {
	http.Handle("/", Handler(defaultHandler))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.zionsoft.net", http.StatusNotFound)
}

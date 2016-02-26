/*
 * Copyright (c) 2016 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package core

import (
	"net/http"
	"reflect"
	"runtime/debug"

	"appengine"
)

type Handler func(http.ResponseWriter, *http.Request)

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			code := http.StatusInternalServerError
			if reflect.TypeOf(e).String() == "*core.Error" {
				code = e.(*Error).Code
			}

			if code == http.StatusInternalServerError {
				c := appengine.NewContext(r)
				c.Errorf("Error: %s\nStack trace:\n%s", e.(*Error).Message, debug.Stack())
			}
			w.WriteHeader(code)
		}
	}()

	fn(w, r)
}

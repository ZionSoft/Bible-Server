/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"src/account"
	"src/core"

	"appengine"
	"appengine/urlfetch"
)

func sendNotificationForVerseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		panic(&core.Error{http.StatusMethodNotAllowed, ""})
	}

	// parses the request
	authorizationKey := r.FormValue("authorizationKey")
	book, err := strconv.ParseInt(r.FormValue("book"), 10, 64)
	chapter, err := strconv.ParseInt(r.FormValue("chapter"), 10, 64)
	verse, err := strconv.ParseInt(r.FormValue("verse"), 10, 64)
	utcOffset, err := strconv.ParseInt(r.FormValue("utcOffset"), 10, 64)
	if len(authorizationKey) == 0 || err != nil {
		panic(&core.Error{http.StatusBadRequest, ""})
	}

	// fetched registration IDs
	c := appengine.NewContext(r)
	var registrationIDs []string
	if registrationIDs, err = account.LoadPushNotificationIDs(c, utcOffset); err != nil {
		panic(&core.Error{http.StatusInternalServerError, err.Error()})
	}

	// constructs request to send the request
	v := struct {
		Book    int64 `json:"book"`
		Chapter int64 `json:"chapter"`
		Verse   int64 `json:"verse"`
	}{
		Book:    book,
		Chapter: chapter,
		Verse:   verse,
	}
	var marshalled []byte
	if marshalled, err = json.Marshal(&v); err != nil {
		panic(&core.Error{http.StatusInternalServerError, err.Error()})
	}

	req := struct {
		RegistrationIDs []string `json:"registration_ids"`
		Data            struct {
			Type  string `json:"type"`
			Attrs string `json:"attrs"`
		}
	}{
		RegistrationIDs: registrationIDs,
		Data: struct {
			Type  string `json:"type"`
			Attrs string `json:"attrs"`
		}{
			Type:  "verse",
			Attrs: string(marshalled),
		},
	}
	if marshalled, err = json.Marshal(&req); err != nil {
		panic(&core.Error{http.StatusInternalServerError, err.Error()})
	}

	var httpReq *http.Request
	if httpReq, err = http.NewRequest("POST", "https://android.googleapis.com/gcm/send", strings.NewReader(string(marshalled))); err != nil {
		panic(&core.Error{http.StatusInternalServerError, err.Error()})
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "key="+authorizationKey)

	// sends request to Google Cloud Messaging server
	httpClient := urlfetch.Client(c)
	var httpResp *http.Response
	if httpResp, err = httpClient.Do(httpReq); err != nil {
		panic(&core.Error{http.StatusInternalServerError, err.Error()})
	}

	// TODO parses the reponse and deletes the dead device account
	dump, err := httputil.DumpResponse(httpResp, true)
	if err != nil {
		panic(&core.Error{http.StatusInternalServerError, err.Error()})
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string(dump))
}

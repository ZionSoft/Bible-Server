/*
 * Copyright (c) 2016 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"

	"src/core"
)

func registerDeviceAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		panic(&core.Error{http.StatusMethodNotAllowed, ""})
	}

	var da deviceAccount
	if err := da.unmarshalHTTPRequest(r); err != nil {
		panic(&core.Error{http.StatusBadRequest, ""})
	}

	c := appengine.NewContext(r)
	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		if da.AccountID == 0 {
			// it might be a new device, or the request is fired by old clients
			if existing, ok := daCache.Get(da.PushNotificationID); ok {
				// hit in-memory cache, so it's a known device with old client
				da.AccountID = existing.(deviceAccount).AccountID
				da.Created = existing.(deviceAccount).Created
			} else {
				// missed in-memory cache, should fall back to datastore
				q := datastore.NewQuery("DeviceAccount").
					Ancestor(createDeviceAccountAncestorKey(c)).
					Filter("PushNotificationID =", da.PushNotificationID)
				var existing []*deviceAccount
				keys, err := q.GetAll(c, &existing)
				if err != nil {
					return err
				}

				if len(keys) == 0 {
					// from an unknown device, should create a new account
					key, err := datastore.Put(c, createIncompleteDeviceAccountKey(c), &da)
					if err != nil {
						return err
					}
					da.AccountID = key.IntID()
					daCache.Set(da.PushNotificationID, da)
					return nil
				} else {
					// a known device with old client
					da.AccountID = keys[0].IntID()
					da.Created = existing[0].Created
				}
			}

			_, err := datastore.Put(c, createDeviceAccountKey(c, da.AccountID), &da)
			daCache.Set(da.PushNotificationID, da)
			return err
		} else {
			// update existing account (for request fired by new clients)
			key := strconv.FormatInt(da.AccountID, 10)
			if existing, ok := daCache.Get(key); ok {
				// hit in-memory cache
				da.Created = existing.(deviceAccount).Created
			} else {
				// missed in-memory cache, falls back to datastore
				var existing deviceAccount
				if err := datastore.Get(c, createDeviceAccountKey(c, da.AccountID), &existing); err != nil {
					return err
				}
				da.Created = existing.Created
			}

			if _, err := datastore.Put(c, createDeviceAccountKey(c, da.AccountID), &da); err != nil {
				return err
			}
			daCache.Set(key, da)
		}

		return nil
	}, nil)
	if err != nil {
		panic(&core.Error{http.StatusInternalServerError, err.Error()})
	}

	// writes the response
	w.WriteHeader(http.StatusCreated)

	var resp struct {
		AccountID int64 `json:"accountId"`
	}
	resp.AccountID = da.AccountID
	b, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(b))
}

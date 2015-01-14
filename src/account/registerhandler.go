/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package account

import (
    "encoding/json"
    "fmt"
    "net/http"

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
        // TODO should caches accout info

        if da.AccountID == 0 {
            // it might be a new device, or the request is fired by old clients
            q := datastore.NewQuery("DeviceAccount").Ancestor(createDeviceAccountAncestorKey(c)).Filter("PushNotificationID =", da.PushNotificationID)
            count, err := q.Count(c)
            if err != nil {
                return err
            }

            if count == 0 {
                // create new account
                key, err := datastore.Put(c, createIncompleteDeviceAccountKey(c), &da)
                if err != nil {
                    return err
                }
                da.AccountID = key.IntID()
            } else {
                // update exisiting account (for request fired by old clients)
                var existing []*deviceAccount
                keys, err := q.GetAll(c, &existing)
                if err != nil {
                    return err
                }
                da.Created = existing[0].Created
                key, err := datastore.Put(c, createDeviceAccountKey(c, keys[0].IntID()), &da)
                if err != nil {
                    return err
                }
                da.AccountID = key.IntID()
            }
        } else {
            // update existing account (for request fired by new clients)
            var existing deviceAccount
            key := createDeviceAccountKey(c, da.AccountID)
            if err := datastore.Get(c, key, &existing); err != nil {
                return err
            }
            da.Created = existing.Created

            if _, err := datastore.Put(c, key, &da); err != nil {
                return err
            }
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

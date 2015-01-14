/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package account

import (
    "encoding/json"
    "errors"
    "net/http"
    "strings"
    "time"

    "appengine"
    "appengine/datastore"
)

type deviceAccount struct {
    AccountID          int64  `datastore:"-" json:"accountId"`
    AccountType        string `json:"-"`
    PushNotificationID string `json:"pushNotificationId"`
    UTCOffset          int64  `json:"utcOffset"`
    Locale             string `json:"locale"`
    Country            string `json:"-"`
    Created            int64  `datastore:",noindex" json:"-"`
    Modified           int64  `datastore:",noindex" json:"-"`
}

func (da *deviceAccount) unmarshalHTTPRequest(r *http.Request) error {
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(da)
    if err != nil {
        return err
    }
    if len(da.PushNotificationID) == 0 {
        return errors.New("missing push notification ID")
    }
    if len(da.Locale) == 0 {
        return errors.New("missing locale")
    }

    // TODO supports other devices
    da.AccountType = "android"

    // TODO should be provided by the client
    da.Country = strings.ToLower(r.Header.Get("X-AppEngine-Country"))

    now := time.Now().Unix()
    da.Created = now
    da.Modified = now

    return nil
}

func createDeviceAccountAncestorKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "DeviceAccount", "DeviceAccountAncestor", 0, nil)
}

func createIncompleteDeviceAccountKey(c appengine.Context) *datastore.Key {
    return datastore.NewIncompleteKey(c, "DeviceAccount", createDeviceAccountAncestorKey(c))
}

func createDeviceAccountKey(c appengine.Context, id int64) *datastore.Key {
    return datastore.NewKey(c, "DeviceAccount", "", id, createDeviceAccountAncestorKey(c))
}

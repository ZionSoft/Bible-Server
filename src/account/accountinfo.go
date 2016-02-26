/*
 * Copyright (c) 2016 ZionSoft. All rights reserved.
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

	"github.com/ZionSoft/lrucache"
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

var daCache *lrucache.LRUCache

var daCacheSize uint64 = 10 * 1024 * 1024

func initializeDeviceAccountCache() {
	daCache = lrucache.New(daCacheSize)
}

func (da deviceAccount) Size() uint64 {
	return uint64(8 + len(da.AccountType) + len(da.PushNotificationID) + 8 + len(da.Locale) + len(da.Country) + 8 + 8)
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

func LoadPushNotificationIDs(c appengine.Context, utcOffset int64) ([]string, error) {
	q := datastore.NewQuery("DeviceAccount").Ancestor(createDeviceAccountAncestorKey(c)).
		Filter("UTCOffset =", utcOffset).Project("PushNotificationID")
	t := q.Run(c)
	registrationIDs := make([]string, 0)
	var err error
	for {
		var da deviceAccount
		_, err = t.Next(&da)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		registrationIDs = append(registrationIDs, da.PushNotificationID)
	}
	return registrationIDs, nil
}

/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package translation

import (
    "net/http"

    "appengine"
    "appengine/datastore"
    "appengine/memcache"

    "src/core"
)

func loadTranslations(c appengine.Context) []*TranslationInfo {
    var translations []*TranslationInfo
    memcache.Gob.Get(c, "TranslationInfo", &translations)
    if len(translations) == 0 {
        // missed memcache, loads from datastore
        q := datastore.NewQuery("TranslationInfo")
        keys, err := q.GetAll(c, &translations)
        if err != nil {
            panic(&core.Error{http.StatusInternalServerError, err.Error()})
        }
        for i, t := range translations {
            t.UniqueId = keys[i].IntID()
        }

        // updates memcache
        item := &memcache.Item{
            Key:    "TranslationInfo",
            Object: translations,
        }
        memcache.Gob.Set(c, item)
    }
    return translations
}

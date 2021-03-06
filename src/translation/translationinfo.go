/*
 * Copyright (c) 2016 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package translation

import (
	"sync"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

type translationInfo struct {
	UniqueId  int64             `datastore:"-" json:"uniqueId"`
	Name      string            `datastore:",noindex" json:"name"`
	ShortName string            `datastore:",noindex" json:"shortName"`
	Language  string            `json:"language"`
	BlobKey   appengine.BlobKey `datastore:",noindex" json:"blobKey"`
	Size      int64             `datastore:",noindex" json:"size"`
	Created   int64             `json:"created"`
	Modified  int64             `json:"modified"`
}

var translationCache struct {
	mu           sync.Mutex
	translations []*translationInfo
}

func loadTranslations(c appengine.Context, forceRefresh bool) ([]*translationInfo, error) {
	translationCache.mu.Lock()
	defer translationCache.mu.Unlock()

	if !forceRefresh {
		if len(translationCache.translations) > 0 {
			return translationCache.translations, nil
		}

		memcache.Gob.Get(c, "TranslationInfo", &translationCache.translations)
		if len(translationCache.translations) > 0 {
			return translationCache.translations, nil
		}
	}

	var err error
	translationCache.translations, err = loadTranslationsFromDatastore(c)
	if err != nil {
		return nil, err
	}

	// updates memcache
	item := &memcache.Item{
		Key:    "TranslationInfo",
		Object: translationCache.translations,
	}
	memcache.Gob.Set(c, item)

	return translationCache.translations, nil
}

func loadTranslationsFromDatastore(c appengine.Context) ([]*translationInfo, error) {
	var translations []*translationInfo
	q := datastore.NewQuery("TranslationInfo")
	keys, err := q.GetAll(c, &translations)
	if err != nil {
		return nil, err
	}
	for i, t := range translations {
		t.UniqueId = keys[i].IntID()
	}
	return translations, nil
}

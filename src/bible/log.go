/*
 * Copyright (c) 2013 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package bible

import (
    "appengine"
    "appengine/datastore"
)

type TranslationDownloadLog struct {
    BlobKey appengine.BlobKey `datastore:"-" json:"blobKey"`
    Count   int64             `datastore:",noindex" json:"count"`
}

func logTranslationDownload(c appengine.Context, blobKey string) error {
    datastore.RunInTransaction(c, func(c appengine.Context) error {
        key := datastore.NewKey(c, "TranslationDownloadLog", blobKey, 0, nil)
        var log TranslationDownloadLog
        if err := datastore.Get(c, key, &log); err != nil && err != datastore.ErrNoSuchEntity {
            return err
        }
        log.Count++
        if _, err := datastore.Put(c, key, &log); err != nil {
            return err
        }
        return nil
    }, nil)
    return nil
}

/*
 * Copyright (c) 2014 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package bible

import (
    "appengine"
)

type TranslationInfo struct {
    UniqueId  int64             `datastore:"-" json:"uniqueId"`
    Name      string            `datastore:",noindex" json:"name"`
    ShortName string            `datastore:",noindex" json:"shortName"`
    Language  string            `json:"language"`
    BlobKey   appengine.BlobKey `datastore:",noindex" json:"blobKey"`
    Size      int64             `datastore:",noindex" json:"size"`
    Timestamp int64             `json:"timestamp"`
}

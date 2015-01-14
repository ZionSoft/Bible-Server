/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package translation

import (
    "appengine"
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

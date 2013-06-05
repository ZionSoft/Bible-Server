package bible

import (
    "appengine"
)

type TranslationInfo struct {
    UniqueId  int64             `datastore:"-"`
    Name      string            `datastore:",noindex" json:"name"`
    ShortName string            `datastore:",noindex" json:"shortName"`
    Language  string            `json:"language"`
    BlobKey   appengine.BlobKey `datastore:",noindex" json:"blobKey"`
    Size      int64             `datastore:",noindex" json:"size"`
    Timestamp int64             `json:"timestamp"`
}

package bible

import (
    "archive/zip"
    "encoding/json"
    "html/template"
    "net/http"
    "time"

    "appengine"
    "appengine/blobstore"
    "appengine/datastore"
)

var uploadTranslationViewTemplate = template.Must(template.New("uploadTranslationView").Parse(uploadTranslationViewTemplateHTML))

const uploadTranslationViewTemplateHTML = `
<html>
<head>
  <link rel="icon" type="image/x-icon" href="https://zionsoft-bible.appspot.com/view/favicon.ico" />
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
  <title>ZionSoft</title>
</head>
<body>
  <div>Add a new translation</div>
  <form action="{{.}}" method="POST" enctype="multipart/form-data">
    <table>
      <tr><td>File:</td><td><input type="file" name="file"></td></tr>
    </table>
    <input type="submit" name="submit" value="Submit">
  </form>
</body>
</html>
`

func uploadTranslationViewHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    uploadUrl, err := blobstore.UploadURL(c, "/admin/uploadTranslation", nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/html")
    err = uploadTranslationViewTemplate.Execute(w, uploadUrl)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
}

type BooksInfo struct {
    Name      string `json:"name"`
    ShortName string `json:"shortName"`
    Language  string `json:"language"`
}

func uploadTranslationHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    blobs, _, err := blobstore.ParseUpload(r)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    blobInfos := blobs["file"]
    if len(blobInfos) != 1 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    blobInfo := blobInfos[0]

    reader, err := zip.NewReader(blobstore.NewReader(c, blobInfo.BlobKey), blobInfo.Size)
    if err != nil {
        blobstore.Delete(c, blobInfo.BlobKey)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    var booksInfo BooksInfo
    for _, f := range reader.File {
        if f.Name != "books.json" {
            continue
        }

        rc, err := f.Open()
        if err != nil {
            blobstore.Delete(c, blobInfo.BlobKey)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        defer rc.Close()
        b := make([]byte, 4096)
        n, err := rc.Read(b)
        if err != nil {
            blobstore.Delete(c, blobInfo.BlobKey)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        err = json.Unmarshal(b[:n - 1], &booksInfo)
        if err != nil {
            blobstore.Delete(c, blobInfo.BlobKey)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }

    var translationInfo TranslationInfo
    translationInfo.Name = booksInfo.Name
    translationInfo.ShortName = booksInfo.ShortName
    translationInfo.Language = booksInfo.Language
    translationInfo.BlobKey = blobInfo.BlobKey
    translationInfo.Size = blobInfo.Size
    translationInfo.Timestamp = time.Now().Unix()

    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "TranslationInfo", nil), &translationInfo)
    if err != nil {
        blobstore.Delete(c, blobInfo.BlobKey)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}

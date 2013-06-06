package bible

import (
    "archive/zip"
    "encoding/json"
    "fmt"
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

func uploadTranslationViewHandler(w http.ResponseWriter, r *http.Request) *appError {
    c := appengine.NewContext(r)
    uploadUrl, err := blobstore.UploadURL(c, "/admin/uploadTranslation", nil)
    if err != nil {
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationViewHandler: Failed to create upload URL '%s'.", err.Error())}
    }
    w.Header().Set("Content-Type", "text/html")
    if err = uploadTranslationViewTemplate.Execute(w, uploadUrl); err != nil {
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationViewHandler: Failed to parse HTML template '%s'.", err.Error())}
    }
    return nil
}

type BooksInfo struct {
    Name      string `json:"name"`
    ShortName string `json:"shortName"`
    Language  string `json:"language"`
}

func uploadTranslationHandler(w http.ResponseWriter, r *http.Request) *appError {
    c := appengine.NewContext(r)
    blobs, _, err := blobstore.ParseUpload(r)
    if err != nil {
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to parse uploaded blob '%s'.", err.Error())}
    }
    blobInfos := blobs["file"]
    if len(blobInfos) != 1 {
        w.WriteHeader(http.StatusBadRequest)
        return &appError{http.StatusBadRequest, string("uploadTranslationHandler: No files uploaded.")}
    }
    blobInfo := blobInfos[0]

    reader, err := zip.NewReader(blobstore.NewReader(c, blobInfo.BlobKey), blobInfo.Size)
    if err != nil {
        blobstore.Delete(c, blobInfo.BlobKey)
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to create blob reader '%s'.", err.Error())}
    }

    var booksInfo BooksInfo
    for _, f := range reader.File {
        if f.Name != "books.json" {
            continue
        }

        rc, err := f.Open()
        if err != nil {
            blobstore.Delete(c, blobInfo.BlobKey)
            return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to open books.json '%s'.", err.Error())}
        }
        defer rc.Close()
        b := make([]byte, 4096)
        n, err := rc.Read(b)
        if err != nil {
            blobstore.Delete(c, blobInfo.BlobKey)
            return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to read books.json '%s'.", err.Error())}
        }

        err = json.Unmarshal(b[:n-1], &booksInfo)
        if err != nil {
            blobstore.Delete(c, blobInfo.BlobKey)
            return &appError{http.StatusBadRequest, fmt.Sprintf("uploadTranslationHandler: Malformed books.json '%s'.", err.Error())}
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
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to save translation info to datastore '%s'.", err.Error())}
    }

    return nil
}

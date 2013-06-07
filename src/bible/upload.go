package bible

import (
    "archive/zip"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "appengine"
    "appengine/blobstore"
    "appengine/datastore"
)

func uploadTranslationHandler(w http.ResponseWriter, r *http.Request) *appError {
    if r.Method != "POST" {
        return &appError{http.StatusMethodNotAllowed, fmt.Sprintf("uploadTranslationHandler: Method '%s' not allowed.", r.Method)}
    }

    // reads the uploaded archive
    f, _, err := r.FormFile("translation")
    if err != nil {
        return &appError{http.StatusBadRequest, string("uploadTranslationHandler: Missing 'translation' file.")}
    }
    defer f.Close()

    buf := make([]byte, r.ContentLength)
    n, err := f.Read(buf)
    if err != nil {
        return &appError{http.StatusInternalServerError, string("uploadTranslationHandler: Failed to read uploaded translation file.")}
    }
    buf = buf[:n]

    // gets part of the translation info from books.json
    var translationInfo TranslationInfo
    translationInfo.Timestamp = time.Now().Unix()

    reader, err := zip.NewReader(f, int64(n))
    if err != nil {
        return &appError{http.StatusBadRequest, string("uploadTranslationHandler: Not a valid zip file.")}
    }
    translationInfo.Size = int64(n)

    for _, file := range reader.File {
        if file.Name != "books.json" {
            continue
        }

        rc, err := file.Open()
        if err != nil {
            return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to open books.json '%s'.", err.Error())}
        }
        defer rc.Close()
        b := make([]byte, 4096) // 4096 bytes should be big enough
        n, err = rc.Read(b)
        if err != nil {
            return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to read books.json '%s'.", err.Error())}
        }
        b = b[:n]

        err = json.Unmarshal(b, &translationInfo)
        if err != nil {
            return &appError{http.StatusBadRequest, fmt.Sprintf("uploadTranslationHandler: Malformed books.json '%s'.", err.Error())}
        }
        break
    }

    // saves uploaded translation into blobstore
    c := appengine.NewContext(r)
    writer, err := blobstore.Create(c, "application/zip")
    if err != nil {
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to create blobstore writer '%s'.", err.Error())}
    }

    _, err = writer.Write(buf)
    if err != nil {
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to write translation file to blobstore '%s'.", err.Error())}
    }
    err = writer.Close()
    if err != nil {
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to flush translation file to blobstore '%s'.", err.Error())}
    }

    // gets rest of the translation info and writes to datastore
    translationInfo.BlobKey, err = writer.Key()
    if err != nil {
        return &appError{http.StatusInternalServerError, string("uploadTranslationHandler: Failed to get blobstore key.")}
    }

    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "TranslationInfo", nil), &translationInfo)
    if err != nil {
        blobstore.Delete(c, translationInfo.BlobKey)
        return &appError{http.StatusInternalServerError, fmt.Sprintf("uploadTranslationHandler: Failed to save translation info to datastore '%s'.", err.Error())}
    }

    return nil
}

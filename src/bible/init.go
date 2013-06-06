package bible

import (
    "net/http"
)

func init() {
    http.Handle("/1.0/downloadTranslation", appHandler(downloadTranslationHandler))
    http.Handle("/1.0/translations", appHandler(queryTranslationsHandler))

    http.Handle("/admin/view/uploadTranslation", appHandler(uploadTranslationViewHandler))
    http.Handle("/admin/uploadTranslation", appHandler(uploadTranslationHandler))
}

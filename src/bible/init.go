package bible

import (
    "net/http"
)

func init() {
    http.HandleFunc("/1.0/downloadTranslation", downloadTranslationHandler)
    http.HandleFunc("/1.0/translations", queryTranslationsHandler)

    http.HandleFunc("/admin/view/uploadTranslation", uploadTranslationViewHandler)
    http.HandleFunc("/admin/uploadTranslation", uploadTranslationHandler)
}

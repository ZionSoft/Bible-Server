/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package notification

import (
	"net/http"

	"src/core"
)

func init() {
	http.Handle("/admin/notification/send/verse", core.Handler(sendNotificationForVerseHandler))
}

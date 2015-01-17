/*
 * Copyright (c) 2015 ZionSoft. All rights reserved.
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 */

package account

import (
    "net/http"

    "src/core"
)

func init() {
    initializeDeviceAccountCache()

    http.Handle("/v1/account/device", core.Handler(registerDeviceAccountHandler))
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var headerTmpl = `
package {{.}}

import (
	"encoding/xml"
	"time"

	{{/*range .Imports*/}}
		{{/*.*/}}
	{{/*end*/}}
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type Client interface {
	SetHeader(interface{})
	Call(string, interface{}, interface{}) error
}

`

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var opsTmpl = `
{{range .}}
	{{$portType := .Name | makePublic}}
	type {{$portType}} struct {
		client Client
	}

	func New{{$portType}}(c Client) *{{$portType}} {
		return &{{$portType}}{
			client: c,
		}
	}

	func (service *{{$portType}}) AddHeader(header interface{}) {
		service.client.AddHeader(header)
	}

	// Backwards-compatible function: use AddHeader instead
	func (service *{{$portType}}) SetHeader(header interface{}) {
		service.client.AddHeader(header)
	}

	{{range .Operations}}
		{{$faults := len .Faults}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
		{{$soapAction := findSOAPAction .Name $portType}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}

		{{/*if ne $soapAction ""*/}}
		{{if gt $faults 0}}
		// Error can be either of the following types:
		// {{range .Faults}}
		//   - {{.Name}} {{.Doc}}{{end}}{{end}}
		{{if ne .Doc ""}}/* {{.Doc}} */{{end}}
		func (service *{{$portType}}) {{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) (*{{$responseType}}, error) {
			response := new({{$responseType}})
			err := service.client.Call("{{$soapAction}}", {{if ne $requestType ""}}request{{else}}nil{{end}}, response)
			if err != nil {
				return nil, err
			}

			return response, nil
		}
		{{/*end*/}}
	{{end}}
{{end}}
`

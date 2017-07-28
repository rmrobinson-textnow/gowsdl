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
		service.client.AddSoapHeader(header)
	}

	// Backwards-compatible function: use AddHeader instead
	func (service *{{$portType}}) SetHeader(header interface{}) {
		service.client.AddSoapHeader(header)
	}

	{{range .Operations}}
		{{$soapAction := findSOAPAction .Name $portType}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}
		{{$faultType := findType .Fault.Message | replaceReservedWords | makePublic}}

		{{/*if ne $soapAction ""*/}}
		func (service *{{$portType}}) {{makePublic .Name | replaceReservedWords}} (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) (*{{$responseType}}, *{{$faultType}}, error) {
			response := new({{$responseType}})
			fault := new({{$faultType}})
			err := service.client.Do(ctx, "{{$soapAction}}", {{if ne $requestType ""}}request{{else}}nil{{end}}, response, fault)
			if err != nil {
				return nil, nil, err
			}

			return response, fault, nil
		}
		{{/*end*/}}
	{{end}}
{{end}}
`

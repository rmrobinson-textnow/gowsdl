// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var opsTmpl = `
{{range .}}
	{{$portType := .Name | makePublic}}
	type {{$portType}} struct {
		client Client
		url string
	}

	func New{{$portType}}(c Client, url string) *{{$portType}} {
		return &{{$portType}}{
			client: c,
			url: url,
		}
	}

	{{range .Operations}}
		{{$soapAction := findSOAPAction .Name $portType}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}
		{{$faultType := findFaultType .Fault.Message | replaceReservedWords | makePublic}}
		{{$operationTypeName := makePublic .Name | replaceReservedWords}}

		{{/*if ne $soapAction ""*/}}
		type {{$operationTypeName}}Call struct {
			service *{{$portType}}
			Request *soap.Request
			Response *soap.Response

			Action string
			requestData *{{$requestType}}
			ResponseData *{{$responseType}}
			FaultData *{{$faultType}}
		}

		func (service *{{$portType}}) New{{$operationTypeName}}Call(reqData *{{$requestType}}) *{{$operationTypeName}}Call {
			call := &{{$operationTypeName}}Call{
				service: service,
				requestData: reqData,
				Action: "{{$soapAction}}",
				ResponseData: &{{$responseType}}{},
				FaultData: &{{$faultType}}{},
			}

			call.Request = soap.NewRequest(call.Action, service.url, call.requestData, call.ResponseData, call.FaultData)
			return call
		}

		func (c *{{$operationTypeName}}Call) Do(ctx context.Context) error {
			var err error
			c.Response, err = c.service.client.Do(ctx, c.Request)
			if err != nil {
				return err
			}
			return nil
		}
		{{/*end*/}}
	{{end}}
{{end}}
`

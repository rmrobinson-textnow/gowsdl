// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `
{{define "SimpleType"}}
	{{$type := replaceReservedWords .Name | makePublic}}
	{{$exists := $type | checkType}}

	{{if not $exists}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{$baseType := toGoType .Restriction.Base}}
		type {{$type}} {{$baseType}}
		const (
			{{with .Restriction}}
				{{range .Enumeration}}
					{{if .Doc}} {{.Doc | comment}} {{end}}

					{{if eq $baseType "string"}}
						{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{goString .Value}}"
					{{else}}
						{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = {{goString .Value}}
					{{end}}
				{{end}}
			{{end}}
		)

		{{$ignore := $type | addType}}
	{{end}}
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base}}
	{{if $baseType}}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{if not .Type}}
			{{ .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base}} ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
		{{else}}
			{{ .Name | makeFieldPublic}} {{toGoType .Type}} ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
		{{end}}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base}}{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeInline"}}
	{{if .SimpleType}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
	{{else}}
		{{replaceReservedWords .Name | makePublic}} struct {
		{{with .ComplexType}}
			{{if ne .ComplexContent.Extension.Base ""}}
				{{template "ComplexContent" .ComplexContent}}
			{{else if ne .SimpleContent.Extension.Base ""}}
				{{template "SimpleContent" .SimpleContent}}
			{{else}}
				{{template "Elements" .Sequence}}
				{{template "Elements" .Choice}}
				{{template "Elements" .SequenceChoice}}
				{{template "Elements" .All}}
				{{template "Attributes" .Attributes}}
			{{end}}
		{{end}}
		} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
	{{end}}
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Ref | toGoType}} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
			{{if not .Type}}
				{{template "ComplexTypeInline" .}}
			{{else}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{replaceReservedWords .Name | makeFieldPublic}} {{if and (.MaxOccurs) (ne .MaxOccurs "1")}}[]{{end}}{{.Type | toGoType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
			{{end}}
		{{end}}
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := .TargetNamespace }}

	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}

	{{range .Elements}}
		{{if not .Type}}
			{{/* ComplexTypeLocal */}}
			{{$name := .Name}}
			{{with .ComplexType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				type {{$name | replaceReservedWords | makePublic}} struct {
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `

					{{if ne .ComplexContent.Extension.Base ""}}
						{{template "ComplexContent" .ComplexContent}}
					{{else if ne .SimpleContent.Extension.Base ""}}
						{{template "SimpleContent" .SimpleContent}}
					{{else}}
						{{template "Elements" .Sequence}}
						{{template "Elements" .Choice}}
						{{template "Elements" .SequenceChoice}}
						{{template "Elements" .All}}
						{{template "Attributes" .Attributes}}
					{{end}}
				}
			{{end}}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$name := replaceReservedWords .Name | makePublic}}
		{{$exists := $name | checkType}}

		{{if not $exists}}
			{{if .Doc}} {{.Doc | comment}} {{end}}
			type {{$name}} struct {
				{{$typeName := .Name | lookupElementType}}
				{{if (and (ne $typeName "") (ne $typeName .Name))}}
					// Element mapping found, will serialize using element name {{$typeName}} instead of {{.Name}}
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$typeName}}\"`" + `
				{{else}}
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{.Name}}\"`" + `
				{{end}}

				{{if ne .ComplexContent.Extension.Base ""}}
					{{template "ComplexContent" .ComplexContent}}
				{{else if ne .SimpleContent.Extension.Base ""}}
					{{template "SimpleContent" .SimpleContent}}
				{{else}}
					{{template "Elements" .Sequence}}
					{{template "Elements" .Choice}}
					{{template "Elements" .SequenceChoice}}
					{{template "Elements" .All}}
					{{template "Attributes" .Attributes}}
				{{end}}
			}
			{{$ignore := $name | addType}}
		{{end}}
	{{end}}
{{end}}
`

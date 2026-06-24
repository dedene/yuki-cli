package api

import (
	"html"
	"strings"
)

func Envelope(operation string, params []Param) string {
	var b strings.Builder
	b.WriteString(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:they="`)
	b.WriteString(Namespace)
	b.WriteString(`">`)
	b.WriteString(`<soapenv:Header/>`)
	b.WriteString(`<soapenv:Body>`)
	b.WriteString(`<they:`)
	b.WriteString(operation)
	b.WriteString(`>`)
	for _, param := range params {
		b.WriteString(`<they:`)
		b.WriteString(param.Name)
		b.WriteString(`>`)
		b.WriteString(html.EscapeString(param.Value))
		b.WriteString(`</they:`)
		b.WriteString(param.Name)
		b.WriteString(`>`)
	}
	b.WriteString(`</they:`)
	b.WriteString(operation)
	b.WriteString(`>`)
	b.WriteString(`</soapenv:Body>`)
	b.WriteString(`</soapenv:Envelope>`)
	return b.String()
}

func SOAPAction(operation string) string {
	return Namespace + operation
}

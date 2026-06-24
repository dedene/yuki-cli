package api

import (
	"strings"
	"testing"
)

func TestEnvelopeEscapesParametersAndUsesYukiNamespace(t *testing.T) {
	got := Envelope("Authenticate", []Param{
		{Name: "accessKey", Value: `abc<&>"'`},
	})

	wantParts := []string{
		`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:they="http://www.theyukicompany.com/">`,
		`<they:Authenticate>`,
		`<they:accessKey>abc&lt;&amp;&gt;&#34;&#39;</they:accessKey>`,
		`</they:Authenticate>`,
	}
	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Fatalf("Envelope() missing %q in:\n%s", part, got)
		}
	}
}

func TestSOAPActionUsesYukiNamespace(t *testing.T) {
	got := SOAPAction("Domains")
	want := "http://www.theyukicompany.com/Domains"
	if got != want {
		t.Fatalf("SOAPAction() = %q, want %q", got, want)
	}
}

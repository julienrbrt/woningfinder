package middleware_test

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strings"

	"github.com/woningfinder/woningfinder/pkg/networking"
)

type lastRequest struct {
	req  *networking.Request
	body []byte
}

type mockClient struct {
	reverseTraceHooks bool
	dnsLookupErr      error
	connectErr        error
	tlsHandshakeErr   error
	err               error
	lastRequests      []lastRequest
}

func (c *mockClient) Send(r *networking.Request) (*networking.Response, error) {
	lastRequest := lastRequest{
		req: r,
	}
	if r.Body != nil {
		b, _ := ioutil.ReadAll(r.Body)
		lastRequest.body = b
	}
	c.lastRequests = append(c.lastRequests, lastRequest)
	trace := httptrace.ContextClientTrace(r.GetContext())
	if trace != nil {
		if !c.reverseTraceHooks {
			trace.DNSStart(httptrace.DNSStartInfo{})
		}
		trace.DNSDone(httptrace.DNSDoneInfo{Err: c.dnsLookupErr})
		if c.reverseTraceHooks {
			trace.DNSStart(httptrace.DNSStartInfo{})
		}
		if c.dnsLookupErr != nil {
			return nil, c.dnsLookupErr
		}

		if !c.reverseTraceHooks {
			trace.ConnectStart("", "")
		}
		trace.ConnectDone("", "", c.connectErr)
		if c.reverseTraceHooks {
			trace.ConnectStart("", "")
		}
		if c.connectErr != nil {
			return nil, c.connectErr
		}

		if !c.reverseTraceHooks {
			trace.TLSHandshakeStart()
		}
		trace.TLSHandshakeDone(tls.ConnectionState{}, c.tlsHandshakeErr)
		if c.reverseTraceHooks {
			trace.TLSHandshakeStart()
		}
		if c.tlsHandshakeErr != nil {
			return nil, c.tlsHandshakeErr
		}

		if !c.reverseTraceHooks {
			trace.WroteRequest(httptrace.WroteRequestInfo{})
		}
		trace.GotFirstResponseByte()
		if c.reverseTraceHooks {
			trace.WroteRequest(httptrace.WroteRequestInfo{})
		}
	}

	if c.err != nil {
		return nil, c.err
	}
	return &networking.Response{
		Request:    r,
		StatusCode: 200,
		Headers:    make(map[string][]string),
		Body:       ioutil.NopCloser(strings.NewReader("Response body")),
		RawResponse: &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("Response body")),
		},
	}, nil
}

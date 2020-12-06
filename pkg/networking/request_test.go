package networking

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestRequest_GetMethodDefault(t *testing.T) {
	r := Request{}
	if r.GetMethod() != http.MethodGet {
		t.Errorf("Expected method to be GET, got %s", r.GetMethod())
	}
}

func TestRequest_GetMethod(t *testing.T) {
	r := Request{Method: http.MethodHead}
	if r.GetMethod() != http.MethodHead {
		t.Errorf("Expected method to be HEAD, got %s", r.GetMethod())
	}
}

func TestRequest_CopyBody(t *testing.T) {
	r := Request{
		Body: strings.NewReader("request-content"),
	}

	body, err := r.CopyBody()
	if err != nil {
		t.Fatalf("Expected CopyBody not to return an error, got %s", err)
	}
	if string(body) != "request-content" {
		t.Errorf("Expected the request body to be 'request-content', got %s", string(body))
	}

	originalBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Expected read all from original request body reader not to return an error, got %s", err)
	}
	if string(originalBody) != "request-content" {
		t.Errorf("Expected the original request body to be 'request-content', got %s", string(originalBody))
	}
}

func TestRequest_CopyBody_BytesBuffer(t *testing.T) {
	r := Request{
		Body: bytes.NewBuffer([]byte("request-content")),
	}

	body, err := r.CopyBody()
	if err != nil {
		t.Fatalf("Expected CopyBody not to return an error, got %s", err)
	}
	if string(body) != "request-content" {
		t.Errorf("Expected the request body to be 'request-content', got %s", string(body))
	}

	originalBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Expected read all from original request body reader not to return an error, got %s", err)
	}
	if string(originalBody) != "request-content" {
		t.Errorf("Expected the original request body to be 'request-content', got %s", string(originalBody))
	}
}

func TestRequest_CopyBodyNil(t *testing.T) {
	r := Request{}
	body, err := r.CopyBody()
	if err != nil {
		t.Fatalf("Expected no error form CopyBody when the body is nil, got %s", err)
	}
	if body != nil {
		t.Errorf("Expected the response body to be nil, got %q", string(body))
	}
}

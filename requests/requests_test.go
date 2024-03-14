package requests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := newClient("http://example.com", http.MethodGet, nil)
	assert.NotNil(t, client)
	assert.Equal(t, "http://example.com", client.url)
	assert.Equal(t, http.MethodGet, client.method)
	assert.NotNil(t, client.header)
	assert.NotNil(t, client.params)
	assert.NotNil(t, client.form)
}

func TestGet(t *testing.T) {
	client := Get("http://example.com")
	assert.NotNil(t, client)
	assert.Equal(t, "http://example.com", client.url)
	assert.Equal(t, http.MethodGet, client.method)
}

func TestPost(t *testing.T) {
	client := Post("http://example.com")
	assert.NotNil(t, client)
	assert.Equal(t, "http://example.com", client.url)
	assert.Equal(t, http.MethodPost, client.method)
}

// Add other HTTP method tests (Put, Delete, Request) similarly

func TestParams(t *testing.T) {
	client := newClient("http://example.com", http.MethodGet, nil)
	client.Params(url.Values{"key": {"value"}})
	assert.Equal(t, "value", client.params.Get("key"))
}

func TestHeader(t *testing.T) {
	client := newClient("http://example.com", http.MethodGet, nil)
	client.Header("Content-Type", "application/json")
	assert.Equal(t, "application/json", client.header.Get("Content-Type"))
}

func TestHeaders(t *testing.T) {
	client := newClient("http://example.com", http.MethodGet, nil)
	headers := http.Header{"Content-Type": []string{"application/json"}}
	client.Headers(headers)
	assert.Equal(t, "application/json", client.header.Get("Content-Type"))
}

// Add other Client methods tests (Form, Json, Multipart) similarly

func TestSend_MockRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/", r.URL.String())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newClient(server.URL, http.MethodGet, nil)
	result := client.Send()

	assert.NotNil(t, result)
	assert.Nil(t, result.Err)
	assert.NotNil(t, result.Resp)
}

// Add other Send method tests (StatusOk, Status2xx, Raw, Text, Json, Save) similarly

func TestSend_EmptyBodyRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/", r.URL.String())
		assert.Equal(t, int64(0), r.ContentLength)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newClient(server.URL, http.MethodGet, nil)
	result := client.Send()

	assert.NotNil(t, result)
	assert.Nil(t, result.Err)
	assert.NotNil(t, result.Resp)
}

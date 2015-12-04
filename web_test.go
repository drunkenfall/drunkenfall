package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testServer() *httptest.Server {
	s := NewServer(nil)
	r := s.BuildRouter()
	return httptest.NewServer(r)
}

func TestServeIndexHtml(t *testing.T) {
	assert := assert.New(t)
	s := testServer()
	defer s.Close()

	res, err := http.Get(s.URL)
	html, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(err)
	assert.True(strings.HasPrefix(string(html), "<!doctype html>"))
}

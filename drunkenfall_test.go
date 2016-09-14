package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thiderman/drunkenfall/websockets"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testServer() *httptest.Server {
	// Server tests use the fake production data.
	s := MockServer("production_data.db")
	ws := websockets.NewServer()
	r := s.BuildRouter(ws)
	return httptest.NewServer(r)
}

func TestServeTournaments(t *testing.T) {
	assert := assert.New(t)
	s := testServer()
	defer s.Close()

	res, err := http.Get(s.URL + "/api/towerfall/tournament/")
	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusOK)

	j, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	lt := &TournamentList{}
	json.Unmarshal(j, lt)
	assert.Equal(2, len(lt.Tournaments))

	res.Body.Close()
}

// TestPermissionsDenied checks all URLs that require producer powers
func TestPermissionsDenied(t *testing.T) {
	assert := assert.New(t)
	s := testServer()
	defer s.Close()

	base := "/api/towerfall/"
	urls := []string{
		"/new/",
		"/%s/start/",
		"/%s/next/",
		"/tournament/%s/tryout/0/toggle/",
		"/tournament/%s/tryout/0/commit/",
	}

	for _, url := range urls {
		if strings.Contains(url, "%s") {
			url = fmt.Sprintf(url, "wrestling")
		}
		res, err := http.Get(s.URL + base + url)
		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusUnauthorized)
	}
}

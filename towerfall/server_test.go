package towerfall

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/drunkenfall/drunkenfall/websockets"
	"github.com/stretchr/testify/assert"
)

func testServer() *httptest.Server {
	// Server tests use the fake production data.
	s := MockServer("production_data.db")

	SetupFakeTournament(nil, s, &NewRequest{"a", "a", time.Now(), true})
	SetupFakeTournament(nil, s, &NewRequest{"b", "b", time.Now(), true})
	s.DB.LoadTournaments()

	ws := websockets.NewServer()
	r := s.BuildRouter(ws)
	return httptest.NewServer(r)
}

func TestServeTournaments(t *testing.T) {
	assert := assert.New(t)
	s := testServer()
	defer s.Close()

	res, err := http.Get(s.URL + "/api/tournament/")
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

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

	base := "/api/"
	urls := []string{
		"/new/",
		"/%s/start/",
		"/%s/next/",
		"/%s/join/",
		"/tournament/%s/0/start/",
		"/tournament/%s/0/end/",
		"/tournament/%s/0/commit/",
	}

	for _, url := range urls {
		if strings.Contains(url, "%s") {
			url = fmt.Sprintf(url, "wrestling")
		}
		res, err := http.Get(s.URL + base + url)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, res.StatusCode)
	}
}

package towerfall

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/olahol/melody"
	"github.com/stretchr/testify/assert"
)

func testServer() *httptest.Server {
	// Server tests use the fake production data.
	s := MockServer("production_data.db")

	SetupFakeTournament(nil, s, &NewRequest{"a", "a", time.Now(), "cover", true})
	SetupFakeTournament(nil, s, &NewRequest{"b", "b", time.Now(), "cover", true})
	s.DB.LoadTournaments()

	ws := melody.New()
	r := s.BuildRouter(ws)
	return httptest.NewServer(r)
}

func TestServeTournaments(t *testing.T) {
	assert := assert.New(t)
	s := testServer()
	defer s.Close()

	res, err := http.Get(s.URL + "/api/tournaments/")
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	j, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	lt := &TournamentList{}
	json.Unmarshal(j, lt)
	assert.Equal(2, len(lt.Tournaments))

	res.Body.Close()
}

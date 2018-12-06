package towerfall

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockServer returns a Server{} a with clean test Database{}
func MockServer(t *testing.T) (*Server, func()) {
	conf := ParseConfig()
	conf.DbName = "test_drunkenfall"
	conf.Port = 56513

	usedPeople = make([]string, 0)

	db, serverTeardown := testDatabase(t, conf)

	s := NewServer(conf, db)
	db.Server = s

	return s, func() {
		serverTeardown()
	}
}

// func testServer() *httptest.Server {
// 	// Server tests use the fake production data.
// 	s := MockServer(t)

// 	SetupFakeTournament(nil, s, &NewRequest{"a", "a", time.Now(), "cover", true})
// 	SetupFakeTournament(nil, s, &NewRequest{"b", "b", time.Now(), "cover", true})

// 	ws := melody.New()
// 	r := s.BuildRouter(ws)
// 	return httptest.NewServer(r)
// }

// func TestServeTournaments(t *testing.T) {
// 	assert := assert.New(t)
// 	s := testServer()
// 	defer s.Close()

// 	res, err := http.Get(s.URL + "/api/tournaments/")
// 	assert.Nil(err)
// 	assert.Equal(http.StatusOK, res.StatusCode)

// 	j, err := ioutil.ReadAll(res.Body)
// 	assert.Nil(err)

// 	lt := &TournamentList{}
// 	json.Unmarshal(j, lt)
// 	assert.Equal(2, len(lt.Tournaments))

// 	res.Body.Close()
// }

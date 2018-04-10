package towerfall

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drunkenfall/drunkenfall/faking"
	"github.com/drunkenfall/drunkenfall/websockets"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Setup variables for the cookies.
var (
	CookieStoreName = "user-session"
	CookieStoreKey  = []byte("dtf")
	CookieStore     = sessions.NewCookieStore(CookieStoreKey)
)

// Determines whether websocket updates should be sent or not.
// This is set to false by the Autoplay functions since they spam with
// hundreds of updates that are pointless. It is also reset to true
// once the Autoplay is over.
var broadcasting = true

// Server is an abstraction that runs a web interface
type Server struct {
	DB        *Database
	config    *Config
	router    *gin.Engine
	simulator *Simulator
	ws        *websockets.Server
}

// NewRequest is the request to make a new tournament
type NewRequest struct {
	Name      string    `json:"name" binding:"required"`
	ID        string    `json:"id" binding:"required"`
	Cover     string    `json:"cover" binding:"required"`
	Scheduled time.Time `json:"scheduled" binding:"required"`
	Fake      bool      `json:"fake" binding:"required"`
}

// CommitPlayer is one state for a player in a commit message
type CommitPlayer struct {
	Ups    int    `json:"ups"`
	Downs  int    `json:"downs"`
	Shot   bool   `json:"shot"`
	Reason string `json:"reason"`
}

// CommitRequest is a request to commit a match state
type CommitRequest struct {
	State []CommitPlayer `json:"state"`
}

// SettingsPostRequest is a settings update
type SettingsPostRequest struct {
	Name  string `json:"name"`
	Nick  string `json:"nick"`
	Color string `json:"color"`
}

type FakeNameResponse struct {
	Name    string `json:"name"`
	Numeral string `json:"numeral"`
}

func init() {
	// To get line numbers in log output
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

// NewServer instantiates a server with an active database
func NewServer(config *Config, db *Database) *Server {
	s := Server{
		DB:     db,
		config: config,
	}
	s.ws = websockets.NewServer()
	s.router = s.BuildRouter(s.ws)

	// Give the db a reference to the server.
	// Not the cleanest, but y'know... here we are.
	db.Server = &s
	return &s
}

// NewHandler shows the page to create a new tournament
func (s *Server) NewHandler(c *gin.Context) {
	var req NewRequest
	var t *Tournament

	if !HasPermission(c, PermissionProducer) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Cannot create tournament unless producer"})
		return
	}

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot JSON"})
		return
	}

	t, err = NewTournament(req.Name, req.ID, req.Cover, req.Scheduled, c, s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot create tournament"})
		return
	}

	log.Printf("Created tournament %s!", t.Name)

	// s.DB.Tournaments = append(s.DB.Tournaments, t)
	// s.DB.tournamentRef[t.ID] = t

	c.Redirect(http.StatusTemporaryRedirect, t.URL())
}

// TournamentHandler returns the current state of the tournament
func (s *Server) TournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tournament": s.getTournament(c)})
}

// TournamentListHandler returns a list of all tournaments
func (s *Server) TournamentListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tournaments": s.DB.asMap()})
}

// JoinHandler shows the tournament view and handles tournaments
func (s *Server) JoinHandler(c *gin.Context) {
	if !HasPermission(c, PermissionPlayer) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to sign in to join a tournament"})
		return
	}

	tm := s.getTournament(c)
	ps := PersonFromSession(s, c)
	tm.TogglePlayer(ps.ID)
	c.Redirect(http.StatusTemporaryRedirect, tm.URL())
}

// EditHandler shows the tournament view and handles tournaments
func (s *Server) EditHandler(c *gin.Context) {
	if !HasPermission(c, PermissionProducer) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to be very hax to edit a tournament"})
		return
	}

	ps := PersonFromSession(s, c)

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	t, err := LoadTournament(data, s.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.DB.OverwriteTournament(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("%s has edited %s!", ps.Nick, t.ID)
	err = t.Persist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, t.URL())

}

// BackfillSemisHandler shows the tournament view and handles tournaments
func (s *Server) BackfillSemisHandler(c *gin.Context) {
	if !HasPermission(c, PermissionJudge) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to be a judge to backfill"})
		return
	}

	tm := s.getTournament(c)
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	spl := strings.Split(string(data), ",")
	err = tm.BackfillSemis(c, spl)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, tm.URL())
}

// ReshuffleHandler reshuffles the player order of the tournament
func (s *Server) ReshuffleHandler(c *gin.Context) {
	if !HasPermission(c, PermissionProducer) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to be a producer to reshuffle"})
		return
	}

	tm := s.getTournament(c)
	err := tm.Reshuffle(c)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, tm.URL())
}

// CreditsHandler returns the data object needed to display the
// credits roll.
func (s *Server) CreditsHandler(c *gin.Context) {
	tm := s.getTournament(c)
	credits, err := tm.GetCredits()

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, credits)
}

// StartTournamentHandler starts tournaments
func (s *Server) StartTournamentHandler(c *gin.Context) {
	if !HasPermission(c, PermissionCommentator) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Cannot start tournament unless commentator or above"})
		return
	}

	tm := s.getTournament(c)
	err := tm.StartTournament(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, tm.URL())
}

// UsurpTournamentHandler usurps tournaments
func (s *Server) UsurpTournamentHandler(c *gin.Context) {
	if !HasPermission(c, PermissionCommentator) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Cannot usurp tournament unless commentator or above"})
		return
	}

	tm := s.getTournament(c)
	err := tm.UsurpTournament()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, tm.URL())
}

// AutoplayTournamentHandler usurps tournaments
func (s *Server) AutoplayTournamentHandler(c *gin.Context) {
	if !HasPermission(c, PermissionProducer) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Cannot autoplay tournament unless producer or above"})
		return
	}

	tm := s.getTournament(c)
	tm.AutoplaySection()
	c.Redirect(http.StatusTemporaryRedirect, tm.URL())
}

// MatchHandler is the common function for match operations.
func (s *Server) MatchHandler(c *gin.Context) {
	if !HasPermission(c, PermissionJudge) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Cannot modify match unless judge or above"})
		return
	}

	m, err := s.getMatch(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't get match",
			"error":   err.Error(),
		})
		return
	}

	action, found := c.Params.Get("action")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Action not set"})
		return
	}

	switch action {
	case "start":
		err = m.Start(c)
	case "end":
		err = m.End(c)
	case "reset":
		err = m.Reset()
	default:
		err = errors.New("Unknown action")
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Match action failed",
			"action":  action,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"action":  action,
		"message": "Done",
	})
}

// MatchCommitHandler commits a single round of a match
func (s *Server) MatchCommitHandler(c *gin.Context) {
	if !HasPermission(c, PermissionJudge) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Cannot commit match unless judge or above"})
		return
	}

	var req CommitRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't get CommitRequest",
			"error":   err.Error(),
		})
		return
	}

	commit := NewMatchCommit(req)
	m, err := s.getMatch(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't get match",
			"error":   err.Error(),
		})
		return
	}

	m.Commit(commit)
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

// ClearTournamentHandler removes all test tournaments
func (s *Server) ClearTournamentHandler(c *gin.Context) {
	err := s.DB.ClearTestTournaments()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't clear test tournaments",
			"error":   err.Error(),
		})
		return
	}
	s.TournamentListHandler(c)
}

// ToggleHandler manages who's in the tournament or not
func (s *Server) ToggleHandler(c *gin.Context) {
	if !HasPermission(c, PermissionJudge) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to be a manager to change joins"})
		return
	}

	t := s.getTournament(c)

	id, found := c.Params.Get("person")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get person ID"})
		return
	}

	err := t.TogglePlayer(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't toggle player",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

// SetTimeHandler sets the pause time for the next match
func (s *Server) SetTimeHandler(c *gin.Context) {
	if !HasPermission(c, PermissionCommentator) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to be a commentator to change times"})
		return
	}

	t := s.getTournament(c)

	st, found := c.Params.Get("time")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get time"})
		return
	}

	x, err := strconv.Atoi(st)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't parse time"})
		return
	}

	m, err := t.NextMatch()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get next match"})
		return
	}

	// If the match is already started, we need to bail
	if m.IsStarted() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Match already started"})
		return
	}

	m.SetTime(c, x)
	c.Redirect(http.StatusTemporaryRedirect, m.URL())
}

// PeopleHandler returns a list of all the players registered in the app
func (s *Server) PeopleHandler(c *gin.Context) {
	if err := s.DB.LoadPeople(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "databaz is ded"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"people": s.DB.People})
}

// CastersHandler sets casters
func (s *Server) CastersHandler(c *gin.Context) {
	if !HasPermission(c, PermissionJudge) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to be a judge to toggle casters"})
		return
	}

	tm := s.getTournament(c)
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	spl := strings.Split(string(data), ",")

	if len(spl) > 2 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Too many casters"})
		return
	}

	tm.SetCasters(spl)
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

func (s *Server) StatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, NewSnapshot(s))
}

// UserHandler returns data about the current user
func (s *Server) UserHandler(c *gin.Context) {
	if !HasPermission(c, PermissionPlayer) {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	c.JSON(http.StatusOK, PersonFromSession(s, c))
}

// DisableHandler disables or enables players
func (s *Server) DisableHandler(c *gin.Context) {
	if !HasPermission(c, PermissionProducer) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to be a producer to toggle"})
		return
	}

	id, found := c.Params.Get("person")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get person ID"})
		return
	}

	s.DB.DisablePerson(id)
	c.JSON(http.StatusOK, gin.H{"people": s.DB.People})
}

// LogoutHandler logs out the user
func (s *Server) LogoutHandler(c *gin.Context) {
	p := PersonFromSession(s, c)

	log.Printf("User '%s' logging out", p.Nick)
	p.RemoveCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

// SettingsHandler gets the POST from the user with a settings update
func (s *Server) SettingsHandler(c *gin.Context) {
	var req SettingsPostRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot JSON"})
		return
	}
	log.Print(req)

	p := PersonFromSession(s, c)
	if p == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Person not found"})
		return
	}

	p.UpdatePerson(&req)
	s.DB.SavePerson(p)

	_ = p.StoreCookies(c)

	c.JSON(http.StatusOK, gin.H{"person": p})
}

// FakeNameHandler returns a fake name for a tournament
func (s *Server) FakeNameHandler(c *gin.Context) {
	name, numeral := faking.FakeTournamentTitle()
	c.JSON(http.StatusOK, gin.H{
		"name":    name,
		"numeral": numeral,
	})
}

func (s *Server) startSimulator(c *gin.Context) {
	var err error

	// If we don't already have a simulator, make one
	if s.simulator == nil {
		log.Print("Creating new simulator")
		s.simulator, err = NewSimulator(s)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tm := s.getTournament(c)
	err = s.simulator.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.simulator.Start(tm.ID)
	c.JSON(http.StatusOK, gin.H{"running": true})
}

func (s *Server) stopSimulator(c *gin.Context) {
	s.simulator.Stop()
	c.JSON(http.StatusOK, gin.H{"running": false})
}

func (s *Server) RequireJudge() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := PersonFromSession(s, c)
		if p == nil || p.Userlevel < PermissionJudge {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Permissions required"})
			c.Abort()
		}
		c.Next()
	}
}

// BuildRouter sets up the routes
func (s *Server) BuildRouter(ws *websockets.Server) *gin.Engine {
	router := gin.Default()

	router.Use(sessions.Sessions(CookieStoreName, CookieStore))

	index := "js/index.html"
	if _, err := os.Stat(index); !os.IsNotExist(err) {
		router.LoadHTMLFiles(index)
		router.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(404, nil)
				return
			}
			c.HTML(200, "index.html", gin.H{})
		})
	}

	api := router.Group("/api")

	// Public routes that don't require any authentication
	api.GET("/people/", s.PeopleHandler)
	api.GET("/people/stats/", s.StatsHandler)
	api.GET("/user/", s.UserHandler)
	api.GET("/user/logout/", s.LogoutHandler)
	api.GET("/user/settings/", s.SettingsHandler)
	api.GET("/fake/name/", s.FakeNameHandler)
	api.GET("/tournaments/", s.TournamentListHandler)
	// api.GET("/tournaments/:id", s.TournamentHandler)
	api.GET("/facebook/login", s.handleFacebookLogin)
	api.GET("/facebook/oauth2callback", s.handleFacebookCallback)

	// Protected routes - everything past this points requires that you
	// are at least a judge.
	api.Use(s.RequireJudge())

	api.POST("/user/:person/disable", s.DisableHandler)
	api.DELETE("/tournaments/", s.ClearTournamentHandler)

	api.POST("/tournaments/", s.NewHandler)

	t := api.Group("/tournaments/:id")
	t.GET("/autoplay/", s.AutoplayTournamentHandler)
	t.GET("/credits/", s.CreditsHandler)
	t.GET("/join/", s.JoinHandler)
	t.GET("/reshuffle/", s.ReshuffleHandler)
	t.GET("/time/:time", s.SetTimeHandler)
	t.GET("/toggle/:person", s.ToggleHandler)
	t.GET("/usurp/", s.UsurpTournamentHandler)
	t.GET("/start/", s.StartTournamentHandler)

	t.POST("/backfill/", s.BackfillSemisHandler)
	t.POST("/casters/", s.CastersHandler)
	t.POST("/edit/", s.EditHandler)

	m := t.Group("/match/:index")

	m.POST("/", s.MatchCommitHandler)
	m.POST("/:action/", s.MatchHandler)

	api.POST("/simulator/start/:id", s.startSimulator)
	api.POST("/simulator/stop/:id", s.stopSimulator)

	// Add the fallback static serving
	// router.PathPrefix("/static/js").Handler(http.StripPrefix("/static/js", http.FileServer(http.Dir("./js/dist/static/js"))))
	// router.PathPrefix("/static/css").Handler(http.StripPrefix("/static/css", http.FileServer(http.Dir("./js/dist/static/css"))))
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./js/static/"))))

	// Install the websockets
	// api.Handle("/auto-updater", websocket.Handler(ws.OnConnected))

	return router
}

// Serve serves forever
func (s *Server) Serve() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}

// SendWebsocketUpdate sends an update to all listening sockets
func (s *Server) SendWebsocketUpdate(kind string, data interface{}) error {
	if !broadcasting {
		return nil
	}

	// TODO(thiderman): Is it safe to just off this as a goroutine?
	// There is a situation where a certain test (TestLavaOrb) makes the
	// tests hang repeatedly if this is not a goroutine. This is extra
	// weird since hundreds of other messages have been sent before that.
	go s.ws.SendAll(&websockets.Message{
		Type: kind,
		Data: data,
	})
	return nil
}

// DisableWebsocketUpdates... disables websocket updates.
func (s *Server) DisableWebsocketUpdates() {
	log.Print("Disabling websocket broadcast")
	broadcasting = false
}

// EnableWebsocketUpdates... enables websocket updates.
func (s *Server) EnableWebsocketUpdates() {
	log.Print("Enabling websocket broadcast")
	broadcasting = true
}

func (s *Server) getMatch(c *gin.Context) (*Match, error) {
	tm := s.getTournament(c)

	i, found := c.Params.Get("")
	if !found {
		return nil, errors.New("match index not set in params")
	}

	index, err := strconv.Atoi(i)
	if err != nil {
		return nil, err
	}

	return tm.Matches[index], nil
}

func (s *Server) getTournament(c *gin.Context) *Tournament {
	id, found := c.Params.Get("id")
	if !found {
		log.Printf("going for id in URL but none is there")
		return nil
	}
	tm := s.DB.tournamentRef[id]
	return tm
}

// Redirect creates a JSON Redirect
func (s *Server) Redirect(w http.ResponseWriter, url string) {
	data, err := json.Marshal(JSONMessage{
		Redirect: url,
	})
	if err != nil {
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// HasPermission checks that the user is allowed to do an action
func HasPermission(c *gin.Context, lvl int) bool {
	session := sessions.Default(c)
	l := session.Get("userlevel")

	// Nothing set in the session - no permission
	if l == nil {
		return false
	}

	return l.(int) >= lvl
}

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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"go.uber.org/zap"
)

// Setup variables for the cookies.
var (
	CookieStoreName = "user-session"
	CookieStoreKey  = []byte("dtf")
)

// Determines whether websocket updates should be sent or not.
// This is set to false by the Autoplay functions since they spam with
// hundreds of updates that are pointless. It is also reset to true
// once the Autoplay is over.
var broadcasting = true

// Message is the data to send back
type websocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Server is an abstraction that runs a web interface
type Server struct {
	DB        *Database
	config    *Config
	router    *gin.Engine
	log       *zap.Logger
	ws        *melody.Melody
	publisher *Publisher
}

// NewRequest is the request to make a new tournament
type NewRequest struct {
	Name      string    `json:"name" binding:"required"`
	ID        string    `json:"id" binding:"required"`
	Scheduled time.Time `json:"scheduled" binding:"required"`
	Cover     string    `json:"cover"`
	Fake      bool      `json:"fake"`
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
	Name       string `json:"name"`
	Nick       string `json:"nick"`
	Color      string `json:"color"`
	ArcherType int    `json:"archer_type"`
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
	var err error
	s := Server{
		DB:     db,
		config: config,
	}
	s.ws = melody.New()
	s.router = s.BuildRouter(s.ws)

	// Add zap logging
	s.log, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	s.router.Use(ginzap.Ginzap(s.log, time.RFC3339, true))

	// Add the Rabbit publisher
	s.publisher, err = NewPublisher(config)
	if err != nil {
		log.Fatal(err)
	}

	// Give the db a reference to the server.
	// Not the cleanest, but y'know... here we are.
	db.Server = &s
	return &s
}

// NewHandler shows the page to create a new tournament
func (s *Server) NewHandler(c *gin.Context) {
	var req NewRequest
	var t *Tournament

	plog := s.log.With(zap.String("path", c.Request.URL.Path))

	err := c.BindJSON(&req)
	if err != nil {
		plog.Error("Bind failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot JSON"})
		return
	}

	idlog := plog.With(zap.String("id", req.ID))
	t, err = NewTournament(req.Name, req.ID, req.Cover, req.Scheduled, c, s)
	if err != nil {
		idlog.Info("Creation failed")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot create tournament"})
		return
	}

	err = s.SendWebsocketUpdate("tournament", t)
	if err != nil {
		idlog.Info("Sending websocket update failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't publish websocket"})
		return
	}

	idlog.Info("Tournament created", zap.String("name", t.Name))
	c.JSON(http.StatusOK, gin.H{"redirect": t.URL()})
}

// TournamentHandler returns the current state of the tournament
func (s *Server) TournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tournament": s.getTournament(c)})
}

// TournamentListHandler returns a list of all tournaments
func (s *Server) TournamentListHandler(c *gin.Context) {
	ts, err := s.DB.GetTournaments()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"tournaments": ts})
}

// JoinHandler shows the tournament view and handles tournaments
func (s *Server) JoinHandler(c *gin.Context) {
	if !HasPermission(c, PermissionPlayer) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need to sign in to join a tournament"})
		return
	}

	tm := s.getTournament(c)
	ps := PersonFromSession(s, c)
	err := tm.TogglePlayer(ps.PersonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"redirect": tm.URL()})
}

// CreditsHandler returns the data object needed to display the
// credits roll.
func (s *Server) CreditsHandler(c *gin.Context) {
	tm := s.getTournament(c)
	tlog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("tournament", tm.Slug),
	)

	credits, err := tm.GetCredits()

	if err != nil {
		tlog.Info("Could not get credits", zap.Error(err))
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	tlog.Info("Credits grabbed")
	c.JSON(http.StatusOK, credits)
}

// StartTournamentHandler starts tournaments
func (s *Server) StartTournamentHandler(c *gin.Context) {
	tm := s.getTournament(c)
	tlog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("tournament", tm.Slug),
	)

	err := tm.StartTournament(c)
	if err != nil {
		tlog.Info("Could not start tournament", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tlog.Info("Tournament started")
	c.JSON(http.StatusOK, gin.H{"redirect": tm.URL()})
}

// UsurpTournamentHandler usurps tournaments
func (s *Server) UsurpTournamentHandler(c *gin.Context) {
	tm := s.getTournament(c)
	tlog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("tournament", tm.Slug),
	)

	err := tm.UsurpTournament()
	if err != nil {
		tlog.Info("Could not usurp tournament", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tlog.Info("Tournament usurped")
	c.JSON(http.StatusOK, gin.H{"redirect": tm.URL()})
}

// AutoplayTournamentHandler plays a section of the tournament automatically
func (s *Server) AutoplayTournamentHandler(c *gin.Context) {
	tm := s.getTournament(c)

	err := tm.AutoplaySection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"redirect": tm.URL()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirect": tm.URL()})
}

// PlayerSummariesHandler returns the players for a tournament
func (s *Server) PlayerSummariesHandler(c *gin.Context) {
	tm := s.getTournament(c)

	tlog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("tournament", tm.Slug),
	)

	ps, err := s.DB.GetPlayerSummaries(tm)
	if err != nil {
		tlog.Info("Could not get summaries", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not get summaries"})
		return
	}

	if len(ps) == 0 {
		c.JSON(http.StatusOK, gin.H{"player_summaries": ps})
		return
	}

	ids := []string{}
	for _, p := range ps {
		ids = append(ids, p.PersonID)
	}

	p, err := s.DB.GetPeopleById(ids...)
	if err != nil {
		tlog.Info("Could not get people", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not get people"})
		return
	}

	for x, y := range p {
		ps[x].Person = y
	}

	c.JSON(http.StatusOK, gin.H{"player_summaries": ps})
}

// StartPlayHandler sends a message to the game that it is time to
// start doing whatever is next.
func (s *Server) StartPlayHandler(c *gin.Context) {
	s.log.Info("sending start_play")
	err := s.publisher.Publish(gStartPlay, StartPlayMessage{})
	if err != nil {
		s.log.Info("Publishing failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"sent": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sent": true})
}

// MatchHandler is the common function for match operations.
func (s *Server) MatchHandler(c *gin.Context) {
	plog := s.log.With(zap.String("path", c.Request.URL.Path))
	m, err := s.getMatch(c)
	if err != nil {
		plog.Info("Could not get match")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't get match",
			"error":   err.Error(),
		})
		return
	}

	mlog := plog.With(
		zap.String("tournament", m.Tournament.Slug),
		zap.Int("match", m.Index),
	)

	action, found := c.Params.Get("action")
	if !found {
		mlog.Info("Action was not set")
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
		mlog.Error("Action failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Match action failed",
			"action":  action,
			"error":   err.Error(),
		})
		return
	}

	mlog.Info("Action executed", zap.String("action", action))
	c.JSON(http.StatusOK, gin.H{
		"action":  action,
		"message": "Done",
	})
}

// ClearTournamentHandler removes all test tournaments
func (s *Server) ClearTournamentHandler(c *gin.Context) {
	plog := s.log.With(zap.String("path", c.Request.URL.Path))

	err := s.DB.ClearTestTournaments()
	if err != nil {
		plog.Info("Couldn't clear test tournaments")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't clear test tournaments",
			"error":   err.Error(),
		})
		return
	}

	plog.Info("Test tournaments cleared")
	s.TournamentListHandler(c)
}

// ToggleHandler manages who's in the tournament or not
func (s *Server) ToggleHandler(c *gin.Context) {
	t := s.getTournament(c)
	tlog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("tournament", t.Slug),
	)

	id, found := c.Params.Get("person")
	if !found {
		tlog.Info("Couldn't get person ID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get person ID"})
		return
	}
	pslog := tlog.With(zap.String("player", id))

	err := t.TogglePlayer(id)
	if err != nil {
		pslog.Error("Could not toggle player")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't toggle player",
			"error":   err.Error(),
		})
		return
	}

	pslog.Info("Player toggled")
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

// SetTimeHandler sets the pause time for the next match
func (s *Server) SetTimeHandler(c *gin.Context) {
	t := s.getTournament(c)
	tlog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("tournament", t.Slug),
	)

	st, found := c.Params.Get("time")
	if !found {
		tlog.Error("Couldn't get time")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get time"})
		return
	}

	x, err := strconv.Atoi(st)
	if err != nil {
		tlog.Error("Couldn't parse time")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't parse time"})
		return
	}

	m, err := t.NextMatch()
	if err != nil {
		tlog.Error("Couldn't get next match")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get next match"})
		return
	}

	mlog := tlog.With(zap.Int("match", m.Index))

	// If the match is already started, we need to bail
	if m.IsStarted() {
		mlog.Warn("Match already started")
		c.JSON(http.StatusForbidden, gin.H{"message": "Match already started"})
		return
	}

	err = m.SetTime(c, x)
	if err != nil {
		mlog.Error("Couldn't set time", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error setting time"})
		return
	}

	mlog.Info("Time set", zap.Int("minutes", x))
	c.JSON(http.StatusOK, gin.H{"redirect": m.URL()})
}

// PeopleHandler returns a list of all the players registered in the app
func (s *Server) PeopleHandler(c *gin.Context) {
	plog := s.log.With(zap.String("path", c.Request.URL.Path))
	ps, err := s.DB.GetPeople()
	if err != nil {
		plog.Error("couldn't get people", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "no people"})
		return
	}

	plog.Info("Players returned")

	c.JSON(http.StatusOK, gin.H{"people": ps})
}

// CastersHandler sets casters
func (s *Server) CastersHandler(c *gin.Context) {
	tm := s.getTournament(c)
	tlog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("tournament", tm.Slug),
	)

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		tlog.Error("Couldn't read body")
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	spl := strings.Split(string(data), ",")
	clog := tlog.With(zap.String("casters", string(data)))

	if len(spl) > 2 {
		clog.Error("Too many casters set")
		c.JSON(http.StatusForbidden, gin.H{"message": "Too many casters"})
		return
	}

	err = tm.SetCasters(spl)
	if err != nil {
		clog.Error("Setting casters failed", zap.Error(err))
		c.JSON(http.StatusForbidden, gin.H{"message": "Setting casters failed"})
		return
	}

	clog.Info("Casters set")
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

func (s *Server) StatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, NewSnapshot(s))
}

// UserHandler returns data about the current user
func (s *Server) UserHandler(c *gin.Context) {
	plog := s.log.With(zap.String("path", c.Request.URL.Path))
	if !HasPermission(c, PermissionPlayer) {
		plog.Info("Not signed in")
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	c.JSON(http.StatusOK, PersonFromSession(s, c))
}

// DisableHandler disables or enables players
func (s *Server) DisableHandler(c *gin.Context) {
	plog := s.log.With(zap.String("path", c.Request.URL.Path))

	id, found := c.Params.Get("person")
	if !found {
		plog.Error("Couldn't get person ID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Couldn't get person ID"})
		return
	}

	err := s.DB.DisablePerson(id)
	if err != nil {
		plog.Error("couldn't get people", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "disable failed"})
		return
	}

	plog.Info("Person disabled", zap.String("person", id))

	ps, err := s.DB.GetPeople()
	if err != nil {
		plog.Error("couldn't get people", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "no people"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"people": ps})
}

// LogoutHandler logs out the user
func (s *Server) LogoutHandler(c *gin.Context) {
	p := PersonFromSession(s, c)
	plog := s.log.With(
		zap.String("path", c.Request.URL.Path),
		zap.String("user", p.PersonID),
	)

	err := p.RemoveCookies(c)
	if err != nil {
		plog.Error("Removing cookies failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Logout failed (?)"})
		return
	}

	plog.Info("User logged out")
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

// SettingsHandler gets the POST from the user with a settings update
func (s *Server) SettingsHandler(c *gin.Context) {
	plog := s.log.With(zap.String("path", c.Request.URL.Path))
	var req SettingsPostRequest

	err := c.BindJSON(&req)
	if err != nil {
		plog.Info("Couldn't JSON")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot JSON"})
		return
	}

	p := PersonFromSession(s, c)
	if p == nil {
		plog.Error("Person not found")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Person not found"})
		return
	}

	p.UpdatePerson(&req)
	err = s.DB.SavePerson(p)
	if err != nil {
		plog.Error("Saving failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database failed"})
		return
	}

	err = p.StoreCookies(c)
	if err != nil {
		plog.Error("Setting cookies failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cookies failed"})
		return
	}

	plog.Info("Person saved", zap.String("person", p.PersonID))
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

func (s *Server) RequireJudge() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := PersonFromSession(s, c)

		if p == nil {
			// If the request is from someone that's not signed in, we need to
			// signal that.
			s.log.Info(
				"Unauthorized",
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please log in"})
			c.Abort()

		} else if p.Userlevel < PermissionJudge {
			// If not, we need to check if the permissions are enough
			s.log.Info(
				"Permission denied",
				zap.String("path", c.Request.URL.Path),
				zap.String("person", p.PersonID),
				zap.Int("userlevel", p.Userlevel),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Permission denied"})
			c.Abort()
		}
		c.Next()
	}
}

// BuildRouter sets up the routes
func (s *Server) BuildRouter(ws *melody.Melody) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	cookieStore := cookie.NewStore(CookieStoreKey)

	router.Use(sessions.Sessions(CookieStoreName, cookieStore))

	index := "js/dist/index.html"
	if _, err := os.Stat(index); !os.IsNotExist(err) {
		router.LoadHTMLFiles(index)
		router.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(404, nil)
				return
			}
			c.HTML(200, "index.html", gin.H{})
		})
		router.Static("/static", "./js/dist/static")
	}

	api := router.Group("/api")

	// Public routes that don't require any authentication
	api.GET("/people/", s.PeopleHandler)
	api.GET("/people/stats/", s.StatsHandler)
	api.GET("/user/", s.UserHandler)
	api.GET("/user/logout/", s.LogoutHandler)
	api.POST("/user/settings/", s.SettingsHandler)
	api.GET("/fake/name/", s.FakeNameHandler)
	api.GET("/tournaments/", s.TournamentListHandler)
	// api.GET("/tournaments/:id", s.TournamentHandler)
	api.GET("/facebook/login", s.handleFacebookLogin)
	api.GET("/facebook/oauth2callback", s.handleFacebookCallback)

	// Websockets are auth free
	api.GET("/auto-updater", func(c *gin.Context) {
		s.log.Info("Websocket setup")

		err := ws.HandleRequest(c.Writer, c.Request)
		if err != nil {
			s.log.Error("Handling websocket setup failed", zap.Error(err))
		}
	})

	ws.HandleMessage(func(ms *melody.Session, msg []byte) {
		s.log.Info("Websocket message", zap.String("message", string(msg)))
		err := ws.Broadcast(msg)
		if err != nil {
			s.log.Error("Handling websocket message failed", zap.Error(err))
		}
	})

	// Protected routes - everything past this point requires that you
	// are at least a judge.
	api.Use(s.RequireJudge())

	api.POST("/user/disable/:person", s.DisableHandler)
	api.DELETE("/tournaments/", s.ClearTournamentHandler)

	api.POST("/tournaments/", s.NewHandler)

	t := api.Group("/tournaments/:id")
	t.GET("/players/", s.PlayerSummariesHandler)
	t.GET("/autoplay/", s.AutoplayTournamentHandler)
	t.GET("/credits/", s.CreditsHandler)
	t.GET("/join/", s.JoinHandler)
	t.GET("/time/:time", s.SetTimeHandler)
	t.GET("/toggle/:person", s.ToggleHandler)
	t.GET("/usurp/", s.UsurpTournamentHandler)
	t.GET("/start/", s.StartTournamentHandler)

	t.POST("/play/", s.StartPlayHandler)
	t.POST("/casters/", s.CastersHandler)

	m := t.Group("/match/:index")
	m.POST("/:action/", s.MatchHandler)

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
	go func(kind string, data interface{}) {
		msg := websocketMessage{
			Type: kind,
			Data: data,
		}

		out, err := json.Marshal(msg)
		if err != nil {
			s.log.Warn("cannot marshal", zap.Error(err))
			return
		}

		err = s.ws.Broadcast(out)
		if err != nil {
			s.log.Error("Broadcast failed", zap.Error(err))
			return
		}
	}(kind, data)

	return nil
}

// DisableWebsocketUpdates... disables websocket updates.
func (s *Server) DisableWebsocketUpdates() {
	// log.Printf("%+v", s)
	// s.log.Info("Disabling websocket broadcast")
	broadcasting = false
}

// EnableWebsocketUpdates... enables websocket updates.
func (s *Server) EnableWebsocketUpdates() {
	// s.log.Info("Enabling websocket broadcast")
	broadcasting = true
}

func (s *Server) getMatch(c *gin.Context) (*Match, error) {
	tm := s.getTournament(c)

	i, found := c.Params.Get("index")
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
	tm, err := s.DB.GetTournament(id)
	if err != nil {
		log.Printf("couldn't get tournament: %+v", err)
		return nil
	}
	return tm
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

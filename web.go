package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	storeKey = []byte("dtf")
	store    = sessions.NewFilesystemStore("cookies.jar", storeKey)
)

// Server is an abstraction that runs via a web interface
type Server struct {
	DB     *Database
	router http.Handler
	logger http.Handler
}

// JSONMessage defines a message to be returned to the frontend
type JSONMessage struct {
	Message  string `json:"message"`
	Redirect string `json:"redirect"`
}

// UpdateMessage returns an update to the current tournament
type UpdateMessage struct {
	Tournament *Tournament `json:"tournament"`
}

// UpdateMatchMessage returns an update to the current match
type UpdateMatchMessage struct {
	Match *Match `json:"match"`
}

// NewRequest is the request to make a new tournament
type NewRequest struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// JoinRequest is the request to join a tournament
type JoinRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CommitPlayer is one state for a player in a commit message
type CommitPlayer struct {
	Ups    int    `json:"ups"`
	Downs  int    `json:"downs"`
	Shot   bool   `json:"shot"`
	Reason string `json:"reason"`
}

type CommitRequest struct {
	State []CommitPlayer `json:"state"`
}

// NewServer instantiates a server with an active database
func NewServer(db *Database) *Server {
	s := Server{DB: db}
	s.router = s.BuildRouter()

	http.Handle("/", s.router)
	s.logger = handlers.LoggingHandler(os.Stdout, s.router)

	return &s
}

// StartHandler is the entry point to the web GUI
//
// If not authenticated, it allows for registration.
// If authenticated and a tournament is running, show that tournament.
// If authenticated and no tournament is running, show a list of tournaments.
func (s *Server) StartHandler(w http.ResponseWriter, r *http.Request) {
	t := getTemplates("static/tournaments.html")
	data := struct {
		Tournaments []*Tournament
	}{
		s.DB.Tournaments,
	}

	render(t, w, r, data)
}

// NewHandler shows the page to create a new tournament
func (s *Server) NewHandler(w http.ResponseWriter, r *http.Request) {
	var req NewRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Fatal(err)
	}

	t, _ := NewTournament(req.Name, req.ID, s.DB)
	log.Printf("Created tournament %s!", t.Name)

	s.DB.Tournaments = append(s.DB.Tournaments, t)
	s.DB.tournamentRef[t.ID] = t

	s.redirect(w, t.URL())
}

// TournamentHandler returns the current state of the tournament
func (s *Server) TournamentHandler(w http.ResponseWriter, r *http.Request) {
	canJoin := false
	vars := mux.Vars(r)

	tm := s.DB.tournamentRef[vars["id"]]
	session, _ := store.Get(r, tm.Name)
	if name, ok := session.Values["player"]; ok {
		canJoin = tm.CanJoin(name.(string))
	} else {
		canJoin = true
	}

	out := struct {
		Tournament *Tournament
		CanJoin    bool
	}{
		tm,
		canJoin,
	}

	data, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// JoinHandler shows the tournament view and handles tournaments
func (s *Server) JoinHandler(w http.ResponseWriter, r *http.Request) {
	var req JoinRequest
	tm := s.getTournament(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(req)

	name := req.Name
	color := req.Color

	if !tm.CanJoin(name) {
		http.Error(w, "too many players", 500)
		return
	}
	if color == "" {
		http.Error(w, "need a color", 500)
		return
	}

	err = tm.AddPlayer(name, color)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// TODO: This should not be here...
	_ = tm.SetMatchPointers()

	log.Printf("%s has joined %s!", name, tm.Name)
	session, err := store.Get(r, tm.Name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// TODO: Does not work. :/
	session.Values["player"] = name
	session.Save(r, w)

	s.redirect(w, tm.URL())
}

// StartTournamentHandler starts tournaments
func (s *Server) StartTournamentHandler(w http.ResponseWriter, r *http.Request) {
	tm := s.getTournament(r)
	err := tm.StartTournament()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	s.redirect(w, tm.URL())
}

// NextHandler starts tournaments
func (s *Server) NextHandler(w http.ResponseWriter, r *http.Request) {
	tm := s.getTournament(r)
	m, err := tm.NextMatch()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	s.redirect(w, m.URL())
}

// MatchHandler shows the Match view and handles Matches
func (s *Server) MatchHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Match *Match
	}{
		s.getMatch(r),
	}

	t := getTemplates("static/matchcontrol.html", "static/player.html", "static/match.html")
	render(t, w, r, data)
}

// MatchToggleHandler starts and stops matches
func (s *Server) MatchToggleHandler(w http.ResponseWriter, r *http.Request) {
	m := s.getMatch(r)
	if !m.IsStarted() {
		log.Printf("%s started", m.String())
		m.Start()
	} else {
		log.Printf("%s ended", m.String())
		m.End()
	}

	data, err := json.Marshal(UpdateMessage{
		Tournament: m.Tournament,
	})
	if err != nil {
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// MatchCommitHandler commits a single round of a match
func (s *Server) MatchCommitHandler(w http.ResponseWriter, r *http.Request) {
	var req CommitRequest
	// tm := s.getTournament(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(req)

	m := s.getMatch(r)
	states := req.State
	scores := [][]int{
		[]int{states[0].Ups, states[0].Downs},
		[]int{states[1].Ups, states[1].Downs},
		[]int{states[2].Ups, states[2].Downs},
		[]int{states[3].Ups, states[3].Downs},
	}
	shots := []bool{
		states[0].Shot,
		states[1].Shot,
		states[2].Shot,
		states[3].Shot,
	}

	m.Commit(scores, shots)

	data, err := json.Marshal(UpdateMatchMessage{
		Match: m,
	})
	if err != nil {
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)

	return
}

// ActionHandler handles judge requests for player action
func (s *Server) ActionHandler(w http.ResponseWriter, r *http.Request) {
	m := s.getMatch(r)
	if !m.IsOpen() {
		log.Print("Not allowing actions on non-started matches")
		log.Print("Match not started")
		http.Redirect(w, r, m.URL(), 302)
		return
	}

	vars := mux.Vars(r)
	index, _ := strconv.Atoi(vars["player"])
	m.Players[index].Action(vars["action"], vars["dir"])

	data, err := json.Marshal(UpdateMessage{
		Tournament: m.Tournament,
	})
	if err != nil {
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// TournamentListHandler returns a list of all tournaments
func (s *Server) TournamentListHandler(w http.ResponseWriter, r *http.Request) {
	tournaments := s.DB.Tournaments
	// tournaments := sort.Reverse(s.DB.Tournaments)
	data, err := json.Marshal(tournaments)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// BuildRouter sets up the routes
func (s *Server) BuildRouter() http.Handler {
	n := mux.NewRouter()
	r := n.PathPrefix("/api/towerfall").Subrouter()

	r.HandleFunc("/tournament/", s.TournamentListHandler)
	r.HandleFunc("/tournament/{id}/", s.TournamentHandler)
	r.HandleFunc("/new/", s.NewHandler)
	r.HandleFunc("/{id}/start/", s.StartTournamentHandler)
	r.HandleFunc("/{id}/join/", s.JoinHandler)
	r.HandleFunc("/{id}/next/", s.NextHandler)

	m := r.PathPrefix("/tournament/{id}/{kind:(tryout|runnerup|semi|final)}/{index:[0-9]+}").Subrouter()
	m.HandleFunc("/", s.MatchHandler)
	m.HandleFunc("/toggle/", s.MatchToggleHandler)
	m.HandleFunc("/commit/", s.MatchCommitHandler)
	m.HandleFunc("/{player:[0-3]}/{action}/{dir:(up|down)}", s.ActionHandler)

	return n
}

// Serve serves forever
func (s *Server) Serve() error {
	log.Print("Listening on :45000")
	return http.ListenAndServe(":45000", s.logger)
}

// getTemplates gets a template with the context set to `extra`, with index.html backing it.
func getTemplates(items ...string) *template.Template {
	items = append(items, "static/index.html")
	t, err := template.ParseFiles(items...)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func render(t *template.Template, w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := NewDatabase("production.db")
	if err != nil {
		log.Fatal(err)
	}

	err = NewServer(db).Serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) getMatch(r *http.Request) *Match {
	var m *Match
	vars := mux.Vars(r)

	tm := s.DB.tournamentRef[vars["id"]]
	kind := vars["kind"]
	index, _ := strconv.Atoi(vars["index"])

	if kind == "tryout" {
		m = tm.Tryouts[index]
	} else if kind == "semi" {
		m = tm.Semis[index]
	} else if kind == "final" {
		m = tm.Final
	}

	return m
}

func (s *Server) getTournament(r *http.Request) *Tournament {
	vars := mux.Vars(r)
	tm := s.DB.tournamentRef[vars["id"]]
	return tm
}

// redirect creates a JSON redirect
func (s *Server) redirect(w http.ResponseWriter, url string) {
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

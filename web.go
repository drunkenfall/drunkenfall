package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
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
	// If there is a post to this URL, it means we are making a new tournament
	if r.Method == "POST" {
		name := r.PostFormValue("name")
		id := r.PostFormValue("id")
		t, _ := NewTournament(name, id, s.DB)
		log.Printf("Created tournament %s!", t.Name)
		http.Redirect(w, r, "/", 302)
		return
	}

	// Elsewise, show the GUI
	t := getTemplates("static/new.html")
	render(t, w, r, struct{}{})
}

// TournamentHandler shows the tournament view and handles tournaments
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

	data := struct {
		Tournament *Tournament
		CanJoin    bool
	}{
		tm,
		canJoin,
	}

	t := getTemplates("static/tournament.html", "static/player.html", "static/match.html")
	render(t, w, r, data)
}

// JoinHandler shows the tournament view and handles tournaments
func (s *Server) JoinHandler(w http.ResponseWriter, r *http.Request) {
	t := getTemplates("static/join.html")
	tm := s.getTournament(r)
	data := struct {
		Tournament *Tournament
	}{
		tm,
	}

	if r.Method == "POST" {
		name := r.PostFormValue("name")
		err := tm.AddPlayer(name)
		if err != nil {
			// TODO: Flash error message
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("%s has joined %s!", name, tm.Name)
		session, err := store.Get(r, tm.Name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// TODO: Does not work. :/
		session.Values["player"] = name
		session.Save(r, w)

		http.Redirect(w, r, tm.URL(), 302)
		return
	}

	render(t, w, r, data)
}

// StartTournamentHandler starts tournaments
func (s *Server) StartTournamentHandler(w http.ResponseWriter, r *http.Request) {
	tm := s.getTournament(r)
	err := tm.StartTournament()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, tm.URL(), 302)
}

// NextHandler starts tournaments
func (s *Server) NextHandler(w http.ResponseWriter, r *http.Request) {
	tm := s.getTournament(r)
	m, err := tm.NextMatch()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, m.URL(), 302)
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

	http.Redirect(w, r, m.URL(), 302)
}

// ActionHandler handles judge requests for player action
func (s *Server) ActionHandler(w http.ResponseWriter, r *http.Request) {
	m := s.getMatch(r)
	vars := mux.Vars(r)
	index, _ := strconv.Atoi(vars["player"])
	m.Players[index].Action(vars["action"], vars["dir"])

	http.Redirect(w, r, m.URL(), 302)
}

// BuildRouter sets up the routes
func (s *Server) BuildRouter() http.Handler {
	r := mux.NewRouter()

	m := "/{id}/{kind:(tryout|runnerup|semi|final)}/{index:[0-9]+}"
	r.HandleFunc("/", s.StartHandler)
	r.HandleFunc("/new", s.NewHandler)
	r.HandleFunc("/{id}/", s.TournamentHandler)
	r.HandleFunc("/{id}/start", s.StartTournamentHandler)
	r.HandleFunc("/{id}/join", s.JoinHandler)
	r.HandleFunc("/{id}/next", s.NextHandler)
	r.HandleFunc(m, s.MatchHandler)
	r.HandleFunc(m+"/toggle", s.MatchToggleHandler)
	r.HandleFunc(m+"/{player:[0-3]}/{action}/{dir:(up|down)}", s.ActionHandler)

	return r
}

// Serve serves forever
func (s *Server) Serve() error {
	log.Print("Listening on :3420")
	return http.ListenAndServe(":3420", s.logger)
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

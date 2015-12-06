package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	t, err := template.ParseFiles("static/tournaments.html", "static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Tournaments []*Tournament
	}{
		s.DB.Tournaments,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.ExecuteTemplate(w, "base", data)
}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	f, _ := ioutil.ReadFile("static/index.html")
	fmt.Fprint(w, string(f))
}

// TournamentHandler shows the tournament view and handles tournaments
func (s *Server) TournamentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tournaments were had.\n")
}

// MatchHandler shows the Match view and handles Matches
func (s *Server) MatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The battle raged ever on.\n")
}

// ActionHandler handles judge requests for player action
func (s *Server) ActionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Arrows were fired, archers were killed, and shots were had.\n")
}

// BuildRouter sets up the routes
func (s *Server) BuildRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.StartHandler)
	r.HandleFunc("/{id}/", s.TournamentHandler)
	r.HandleFunc("/{id}/{kind:(tryout|runnerup|semi|final)}/{index:[0-9]+}/", s.MatchHandler)
	r.HandleFunc("/{id}/{kind:(tryout|runnerup|semi|final)}/{index:[0-9]+}/{player:[0-3]}", s.ActionHandler)

	return r
}

// Serve serves forever
func (s *Server) Serve() error {
	log.Print("Listening on :3420")
	return http.ListenAndServe(":3420", s.logger)
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

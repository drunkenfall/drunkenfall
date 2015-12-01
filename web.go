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

// StartHandler is the entry point to the web GUI
//
// If not authenticated, it allows for registration.
// If authenticated and a tournament is running, show that tournament.
// If authenticated and no tournament is running, show a list of tournaments.
func StartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	f, _ := ioutil.ReadFile("html/index.html")
	fmt.Fprint(w, string(f))
}

// TournamentHandler shows the tournament view and handles tournaments
func TournamentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tournaments were had.\n")
}

// MatchHandler shows the Match view and handles Matches
func MatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The battle raged ever on.\n")
}

// ActionHandler handles judge requests for player action
func ActionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Arrows were fired, archers were killed, and shots were had.\n")
}

// BuildRouter sets up the routes
func BuildRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", StartHandler)
	r.HandleFunc("/{id}/", TournamentHandler)
	r.HandleFunc("/{id}/{kind:(tryout|runnerup|semi|final)}/{index:[0-9]+}/", MatchHandler)
	r.HandleFunc("/{id}/{kind:(tryout|runnerup|semi|final)}/{index:[0-9]+}/{player:[0-3]}", ActionHandler)

	return r
}

func main() {
	r := BuildRouter()
	http.Handle("/", r)

	log.Print("Listening on :3420")
	fmt.Println()

	logged := handlers.LoggingHandler(os.Stdout, r)
	err := http.ListenAndServe(":3420", logged)
	if err != nil {
		log.Fatal(err)
	}
}

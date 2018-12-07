package towerfall

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	wTournament      = "tournament"
	wPlayerSummaries = "player_summaries"
	wRunnerups       = "runnerups"
	wPlayer          = "player"
	wMatch           = "match"
	wMatches         = "matches"
	wMatchEnd        = "match_end"
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

type wsTournament struct {
	TournamentID uint        `json:"tournament_id"`
	Tournament   *Tournament `json:"tournament"`
}

type wsPlayerSummaries struct {
	TournamentID    uint             `json:"tournament_id"`
	PlayerSummaries []*PlayerSummary `json:"player_summaries"`
}

type wsRunnerups struct {
	TournamentID uint             `json:"tournament_id"`
	Runnerups    []*PlayerSummary `json:"runnerups"`
}

type wsMatches struct {
	TournamentID uint     `json:"tournament_id"`
	Matches      []*Match `json:"matches"`
}

type wsMatchEnd struct {
	TournamentID    uint             `json:"tournament_id"`
	Tournament      *Tournament      `json:"tournament"`
	Runnerups       []*PlayerSummary `json:"runnerups"`
	PlayerSummaries []*PlayerSummary `json:"player_summaries"`
	Matches         []*Match         `json:"matches"`
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

		s.log.Info(
			"Websocket update",
			zap.String("type", kind),
			zap.Any("data", data),
		)
	}(kind, data)

	return nil
}

// SendTournamentUpdate sends an update about the tournament
func (s *Server) SendTournamentUpdate(t *Tournament) error {
	return s.SendWebsocketUpdate(wTournament, wsTournament{
		TournamentID: t.ID,
		Tournament:   t,
	})
}

// SendPlayerSummariesUpdate sends an update to the player summaries
func (s *Server) SendPlayerSummariesUpdate(t *Tournament) error {
	summaries, err := s.DB.GetPlayerSummaries(t)
	if err != nil {
		return err
	}

	return s.SendWebsocketUpdate(wPlayerSummaries, wsPlayerSummaries{
		TournamentID:    t.ID,
		PlayerSummaries: summaries,
	})
}

// SendRunnerupsUpdate sends an update to the runnerups
func (s *Server) SendRunnerupsUpdate(t *Tournament) error {
	runnerups, err := s.DB.GetAllRunnerups(t)
	if err != nil {
		return err
	}

	return s.SendWebsocketUpdate(wRunnerups, wsRunnerups{
		TournamentID: t.ID,
		Runnerups:    runnerups,
	})
}

// SendPlayerUpdate sends a status update for a single player
func (s *Server) SendPlayerUpdate(m *Match, st *PlayerState) error {
	return s.SendWebsocketUpdate(
		wPlayer,
		PlayerStateUpdateMessage{m.Tournament.ID, m.Index, st},
	)
}

// SendMatchUpdate sends a status update for the entire match
func (s *Server) SendMatchUpdate(m *Match) error {
	// TODO(thiderman): This is terrible
	// Load the match from the database to avoid caching problems
	lm, err := globalDB.GetMatch(m.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	sts, err := globalDB.GetPlayerStates(lm)
	if err != nil {
		return errors.WithStack(err)
	}

	return s.SendWebsocketUpdate(
		wMatch,
		MatchUpdateMessage{
			m.Tournament.ID,
			lm,
			sts,
		},
	)
}

// SendMatchesUpdate sends an update about the matches
func (s *Server) SendMatchesUpdate(t *Tournament) error {
	ms, err := s.DB.GetMatches(t, "all")
	if err != nil {
		return err
	}

	return s.SendWebsocketUpdate(wMatches, wsMatches{
		TournamentID: t.ID,
		Matches:      ms,
	})
}

// SendMatchEndUpdate sends the updates needed when a match is ended
func (s *Server) SendMatchEndUpdate(t *Tournament) error {
	var err error
	s.log.Info("Sending match end update", zap.Uint("tid", t.ID))

	ms, err := s.DB.GetMatches(t, "all")
	if err != nil {
		return err
	}

	runnerups, err := s.DB.GetAllRunnerups(t)
	if err != nil {
		return err
	}

	summaries, err := s.DB.GetPlayerSummaries(t)
	if err != nil {
		return err
	}

	return s.SendWebsocketUpdate(wMatchEnd, wsMatchEnd{
		TournamentID:    t.ID,
		Tournament:      t,
		Matches:         ms,
		Runnerups:       runnerups,
		PlayerSummaries: summaries,
	})
}

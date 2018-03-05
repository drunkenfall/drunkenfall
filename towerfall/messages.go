package towerfall // JSONMessage defines a message to be returned to the frontend
type JSONMessage struct {
	Message  string `json:"message"`
	Redirect string `json:"redirect"`
}

// SettingsUpdateResponse is a redirect response with an extra Person field
type SettingsUpdateResponse struct {
	Redirect string  `json:"redirect"`
	Person   *Person `json:"person"`
}

// GeneralRedirect is an explicit permission failure
type GeneralRedirect JSONMessage

// TournamentMessage returns a single tournament
type TournamentMessage struct {
	Tournament *Tournament `json:"tournament"`
}

// UpdateMessage returns an update to the current tournament
type UpdateMessage TournamentMessage

// UpdateMatchMessage returns an update to the current match
type UpdateMatchMessage struct {
	Match *Match `json:"match"`
}

// TournamentList returns a list with tournaments
type TournamentList struct {
	Tournaments map[string]*Tournament `json:"tournaments"`
}

// UpdateStateMessage returns an update to the current match
type UpdateStateMessage TournamentList

// PeopleList returns a list with users
type PeopleList struct {
	People []*Person `json:"people"`
}

type PlayerStateUpdateMessage struct {
	Tournament string      `json:"tournament"`
	Match      int         `json:"match"`
	Player     int         `json:"player"`
	State      PlayerState `json:"state"`
}

type MatchUpdateMessage struct {
	Tournament string `json:"tournament"`
	Match      *Match `json:"state"`
}

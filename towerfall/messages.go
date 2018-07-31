package towerfall

import "time"

// JSONMessage defines a message to be returned to the frontend
type JSONMessage struct {
	Message  string `json:"message"`
	Redirect string `json:"redirect"`
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

// Constant strings for use as kinds when communicating with the game
const (
	gMatch      = "match"
	gConnect    = "game_connected"
	gDisconnect = "game_disconnected"
)

// GameMatchMessage is the message sent to the game about the
// configuration of the next match
type GameMatchMessage struct {
	Players    []GamePlayer `json:"players"`
	Level      string       `json:"level"`
	Tournament string       `json:"tournament"`
}

// GamePlayer is a player object to be consumed by the game
type GamePlayer struct {
	Name  string `json:"name"`
	Color int    `json:"color"`
}

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

const EnvironmentKill = -1

// Kill reasons
const (
	rArrow = iota
	rExplosion
	rBrambles
	rJumpedOn
	rLava
	rShock
	rSpikeBall
	rFallingObject
	rSquish
	rCurse
	rMiasma
	rEnemy
	rChalice
)

// Arrows
const (
	aNormal = iota
	aBomb
	aSuperBomb
	aLaser
	aBramble
	aDrill
	aBolt
	aToy
	aFeather
	aTrigger
	aPrism
)

// Message types
const (
	inKill       = "kill"
	inRoundStart = "round_start"
	inRoundEnd   = "round_end"
	inMatchStart = "match_start"
	inMatchEnd   = "match_end"
	inPickup     = "arrows_collected"
	inShot       = "arrow_shot"
	inShield     = "shield_state"
	inWings      = "wings_state"
	inOrbLava    = "lava_orb_state"
	// TODO(thiderman): Non-player orbs are not implemented
	inOrbSlow   = "slow_orb_state"
	inOrbDark   = "dark_orb_state"
	inOrbScroll = "scroll_orb_state"
)

type KillMessage struct {
	Player int `json:"player"`
	Killer int `json:"killer"`
	Cause  int `json:"cause"`
}

type ArrowMessage struct {
	Player int    `json:"player"`
	Arrows Arrows `json:"arrows"`
}

type ShieldMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

type WingsMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

type SlowOrbMessage struct {
	State bool `json:"state"`
}

type DarkOrbMessage struct {
	State bool `json:"state"`
}

type ScrollOrbMessage struct {
	State bool `json:"state"`
}

type LavaOrbMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

// List of integers where one item is an arrow type as described in
// the arrow types above.
type Arrows []int

type StartRoundMessage struct {
	Arrows []Arrows `json:"arrows"`
}

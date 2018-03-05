package websockets

// Message is the data to send back
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Ping is a simple ping message
type Ping struct {
	P int `json:"p"`
}

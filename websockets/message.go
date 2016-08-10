package websockets

// Message is the data to send back
type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (m *Message) String() string {
	return m.Author + " says " + m.Body
}

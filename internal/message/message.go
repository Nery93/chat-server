package message

type Message struct {
	Type  string   `json:"type"`
	User  string   `json:"user"`
	Text  string   `json:"text"`
	Users []string `json:"users,omitempty"`
}

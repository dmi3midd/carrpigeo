package email

type Email struct {
	ID       string `json:"id"`
	Sender   string `json:"sender"`
	Reciever string `json:"reciever"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	SentAt   string `json:"sent_at"`
}

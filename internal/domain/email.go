package domain

import "time"

type Email struct {
	ID       string    `json:"id" db:"id"`
	Sender   string    `json:"sender" db:"sender"`
	Reciever string    `json:"reciever" db:"reciever"`
	Subject  string    `json:"subject" db:"subject"`
	Body     string    `json:"body" db:"body"`
	IsHTML   bool      `json:"is_html" db:"is_html"`
	SentAt   time.Time `json:"sent_at" db:"sent_at"`
}

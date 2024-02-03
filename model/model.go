package model

type Message struct {
	To         string `json:"to"`
	SenderName string `json:"senderName"`
	Message    string `json:"message"`
}

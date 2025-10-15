package message

import "time"

type Message struct {
	ID uint `json:"id"`
	UserId uint `json:"user_id"`
	ChatId uint `json:"chat_id"`
	Content []byte `json:"content"`
	Nonce []byte `json:"nonce"`
	CreatedAt time.Time `json:"created_at"`
}

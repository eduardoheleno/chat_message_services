package message

import (
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *Message) error
	GetPaginated(idChat string, offset int) ([]Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db}
}

func (m *messageRepository) Create(message *Message) error {
	err := m.db.Create(message).Error
	if err != nil {
		return err
	}

	return nil
}

func (m *messageRepository) GetPaginated(idChat string, offset int) ([]Message, error) {
	messages := make([]Message, 0)
	query := m.db.Where("chat_id = ?", idChat).
		Order("created_at desc").
		Offset(offset).
		Limit(20)

	err := query.Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

package data

import (
	"database/sql"
	"time"
)

type ChatMessage struct {
	ID        int
	UserID    int
	Role      string
	Content   string
	CreatedAt time.Time
}

type ChatMessageModel struct {
	DB *sql.DB
}

func (m ChatMessageModel) GetHistoryByUserID(userID int, limit int) ([]ChatMessage, error) {
	query := `
		SELECT id, user_id, role, content, created_at
		FROM chat_messages
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := m.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append([]ChatMessage{msg}, messages...) // reverse order
	}
	return messages, nil
}

func (m ChatMessageModel) Insert(userID int, role, content string) error {
	query := `INSERT INTO chat_messages (user_id, role, content) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(query, userID, role, content)
	return err
}

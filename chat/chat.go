package chat

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/johnshiver/asapp_challenge/utils"
)

type Message struct {
	ID          int64          `json:"id"`
	RecipientID int64          `json:"recipient" db:"recipient_id"`
	SenderID    int64          `json:"sender" db:"sender_id"`
	CreatedAt   utils.JSONTime `json:"timestamp" db:"created_at"`
	Content     Content        `json:"content" db:"message_content"`
}

var (
	ErrNoRecord = errors.New("no matching record found")
)

//go:generate sh -c "mockgen -source=chat.go -package chat -destination chat_mock.go Service"
type Service interface {
	Insert(tx *sqlx.Tx, senderID, recipientID int64, messageContent Content) (int64, time.Time, error)
	GetMessagesByRecipient(tx *sqlx.Tx, senderID, recipientID, start int64, limit int) ([]*Message, error)
}

// Implements UserService
type Manager struct {
}

// Custom type for JSONB
// https://www.alexedwards.net/blog/using-postgresql-jsonb
type Content map[string]interface{}

// Scan override the sqlx scanner
func (c Content) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Content) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

func (c Content) Validate() error {
	contentType, ok := c["type"]
	if !ok {
		return fmt.Errorf("content must have a type")
	}

	switch contentType {
	case "text":
		if text, ok := c["text"]; !ok {
			return fmt.Errorf("text type message must include text field")
		} else {
			_, ok := text.(string)
			if !ok {
				return fmt.Errorf("text field must be a string")
			}
		}
	case "image":
		var imageValidationErrors []string

		if urlRaw, ok := c["url"]; !ok {
			imageValidationErrors = append(imageValidationErrors, "image type message must include url field")
		} else {
			url, ok := urlRaw.(string)
			if !ok {
				imageValidationErrors = append(imageValidationErrors, "url field must be a string")
			}
			if !utils.IsValidUrl(url) {
				imageValidationErrors = append(imageValidationErrors, "url field must be a valid url")
			}
		}

		heightRaw, ok := c["height"]
		if !ok {
			imageValidationErrors = append(imageValidationErrors, "image type message must include height field")
		} else {
			height, ok := heightRaw.(float64)
			if !ok {
				imageValidationErrors = append(imageValidationErrors, "height field must be an int")
				fmt.Println(height)
			}
		}

		width, ok := c["width"]
		if !ok {
			imageValidationErrors = append(imageValidationErrors, "image type message must include width field")
		} else {
			_, ok := width.(float64)
			if !ok {
				imageValidationErrors = append(imageValidationErrors, "width field must be an int")
			}
		}

		if len(imageValidationErrors) > 0 {
			return fmt.Errorf(strings.Join(imageValidationErrors, " | "))
		}
	case "video":
		var videoValidationErrors []string
		if urlRaw, ok := c["url"]; !ok {
			videoValidationErrors = append(videoValidationErrors, "video type message must include url field")
		} else {
			if url, ok := urlRaw.(string); !ok {
				videoValidationErrors = append(videoValidationErrors, "url field must be a string")
			} else if !utils.IsValidUrl(url) {
				videoValidationErrors = append(videoValidationErrors, "url field must be a valid url")
			}
		}

		if sourceRaw, ok := c["source"]; !ok {
			videoValidationErrors = append(videoValidationErrors, "video type message must include source field")
		} else {
			if source, ok := sourceRaw.(string); !ok {
				videoValidationErrors = append(videoValidationErrors, "source field must be a string")
			} else {
				var validSources = map[string]struct{}{
					"youtube": {},
					"vimeo":   {},
				}
				if _, ok := validSources[source]; !ok {
					videoValidationErrors = append(videoValidationErrors, "source field must be one of {youtube, vimeo}")
				}
			}
		}
		if len(videoValidationErrors) > 0 {
			return fmt.Errorf(strings.Join(videoValidationErrors, " | "))
		}

	default:
		return fmt.Errorf("message type must be one of {text, image, video}")
	}

	return nil
}

func (m *Manager) Insert(tx *sqlx.Tx, senderId, recipientId int64, messageContent Content) (int64, time.Time, error) {
	const q = `
	  INSERT INTO chat_messages(sender_id, recipient_id, message_content)
	  VALUES ($1, $2, $3)
	  RETURNING id, created_at
    `
	var (
		rowId     int64
		createdAt time.Time
	)

	messageContentJSON, err := json.Marshal(messageContent)
	if err != nil {
		return -1, time.Time{}, err
	}
	err = tx.QueryRowx(q, senderId, recipientId, messageContentJSON).Scan(&rowId, &createdAt)
	if err != nil {
		return -1, time.Time{}, err
	}
	return rowId, createdAt, nil
}

func (m *Manager) GetMessagesByRecipient(tx *sqlx.Tx, senderID, recipientID, startId int64, limit int) ([]*Message, error) {
	// paging through results in postgres https://use-the-index-luke.com/sql/partial-results/fetch-next-page
	const q = `
	 SELECT id, sender_id, recipient_id, message_content, created_at
	 FROM chat_messages 
	 WHERE id >= $1
	 AND sender_id=$2 AND recipient_id=$3
	 ORDER BY id
	 FETCH FIRST $4 ROWS ONLY
    `
	var messages []*Message
	err := tx.Select(&messages, q, startId, senderID, recipientID, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return messages, nil
}

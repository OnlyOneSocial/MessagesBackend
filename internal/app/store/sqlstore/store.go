package sqlstore

import (
	"database/sql"

	"github.com/katelinlis/MessagesBackend/internal/app/store"
	_ "github.com/lib/pq" //db import
)

//Store ...
type Store struct {
	db                *sql.DB
	messageRepository *MessageRepository
}

//New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

//Message ...
func (s *Store) Message() store.MessageRepository {
	if s.messageRepository != nil {
		return s.messageRepository
	}

	s.messageRepository = &MessageRepository{
		store: s,
	}

	return s.messageRepository
}

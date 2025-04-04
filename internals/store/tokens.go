package store

import (
	"database/sql"
	"time"

	"github.com/mrangel-jr/complete-go/internals/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	// AddToken adds a new token to the store.
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, expiry time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID string, scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	// Implementation for inserting the token into the database
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (user_id, hash, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`
	_, err := t.db.Exec(query, token.UserID, token.Hash, token.Expiry, token.Scope)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID string, scope string) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1 AND scope = $2
	`
	_, err := t.db.Exec(query, userID, scope)
	return err
}

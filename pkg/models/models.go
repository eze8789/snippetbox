package models

import (
	"errors"
	"time"
)

// ErrNoRecord Manage errors when record not found
var ErrNoRecord = errors.New("models: No matching record found")

// ErrInvalidCreds manage errors when the user and/or password is incorrect
var ErrInvalidCreds = errors.New("models: invalid credentials")

// ErrDupEmail manage errors when user email is duplicated
var ErrDupEmail = errors.New("models: email duplicated")

// Snippet is a struct matching the DB model
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User struct for users
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

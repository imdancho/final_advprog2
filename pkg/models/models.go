package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	// tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Snippet struct { //reviews
	ID      int
	User_id int
	Content string
	Created string
}

type Service struct {
	ID      int
	Title   string
	Content string
	Master  string
	Price   int
}

type Appointment struct {
	ID         int
	User_id    int
	Service_id string
	Time       string
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
	Role           string
	C_code         int
	Confirmation   bool
}

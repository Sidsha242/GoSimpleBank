package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

//Contains Payload struct and functions to create and verify JWT tokens

//Represents the data stored in a token [struct-as we need concrete type]
type Payload struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	IssuedAt time.Time  `json:"issued_at"`
	ExpiredAt time.Time  `json:"expired_at"`
}

//Create a new payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom() //
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID: id,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}


//Check if the token is valid by comparing the current time with the expiration time

var ErrExpiredToken = errors.New("token has expired")  //global error

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
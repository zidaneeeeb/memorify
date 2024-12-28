package auth

import (
	"context"
	"time"
)

type Service interface {
	// LoginSocial checks the given email and password with
	// the actual data. It returns a token and the encapsulated
	// data if the login process is success.
	LoginSocial(ctx context.Context, tokenEmail string) (string, TokenData, error)

	// GetUserByID returns a user with the given user ID.
	GetUserByID(ctx context.Context, userID string) (User, error)

	// UpdateUser updates existing user
	// with the given user data.
	//
	// UpdateUser do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdateUser(ctx context.Context, reqUser User) error

	// ValidateToken validates the given token and returns the
	// data encapsulated in the token if the given token is
	// valid.
	ValidateToken(ctx context.Context, token string) (TokenData, error)

	// RefreshToken validates the given token and returns a
	// new token with the same encapsulated data but refreshed.
	//
	// RefreshToken is used to avoid expired token.
	RefreshToken(ctx context.Context, token string) (string, error)
}

type User struct {
	ID         string
	Fullname   string
	Username   string
	Email      string
	Type       Type
	Quota      int
	CreateTime time.Time
	UpdateTime time.Time
}

type TokenData struct {
	UserID   string
	Fullname string
	Username string
	Email    string
	Type     Type
	Quota    int
}

type Type int

const (
	TypeUnknown Type = 0
	TypeFree    Type = 1
	TypePemium  Type = 2
	TypePending Type = 3
)

var (
	TypeName = map[Type]string{
		TypeFree:    "free",
		TypePemium:  "premium",
		TypePending: "pending",
	}
)

func (t Type) String() string {
	return TypeName[t]
}

func (t Type) Value() int {
	return int(t)
}

type GetUserAuthFilter struct {
	Email  string
	UserID string
}

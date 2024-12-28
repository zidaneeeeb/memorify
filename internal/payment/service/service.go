package service

import (
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/content"
	"time"
)

// service implements subject.Service.
type service struct {
	pgStore PGStore
	user    auth.Service
	content content.Service
	timeNow func() time.Time
}

// New creates a new service.
func New(pgStore PGStore, user auth.Service, content content.Service) (*service, error) {
	s := &service{
		pgStore: pgStore,
		user:    user,
		content: content,
		timeNow: time.Now,
	}

	return s, nil
}

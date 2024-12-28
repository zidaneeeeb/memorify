package service

import (
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/template"
	"time"
)

// service implements subject.Service.
type service struct {
	pgStore  PGStore
	user     auth.Service
	template template.Service
	timeNow  func() time.Time
}

// New creates a new service.
func New(pgStore PGStore, user auth.Service, template template.Service) (*service, error) {
	s := &service{
		pgStore:  pgStore,
		user:     user,
		template: template,
		timeNow:  time.Now,
	}

	return s, nil
}

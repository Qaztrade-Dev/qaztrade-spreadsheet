package domain

import "context"

type Manager struct {
	UserID   string
	Email    string
	Fullname string
	Roles    []string
}

type ManagersRepository interface {
	GetMany(ctx context.Context) ([]*Manager, error)
	GetCurrent(ctx context.Context, UserID string) (*Manager, error)
}

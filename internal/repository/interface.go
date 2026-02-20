package repository

import (
	"context"
	"shortener/internal/model"
)

type Repository interface {
	SaveLink(ctx context.Context, link *model.Link) error
	GetLink(ctx context.Context, shortURL string) (*model.Link, error)
	LinkExists(ctx context.Context, shortURL string) (bool, error)
	SaveAnalytics(ctx context.Context, a *model.Analytics) error
	GetAnalytics(ctx context.Context, shortURL string, limit int) ([]model.Analytics, error)
	CountAnalytics(ctx context.Context, shortURL string) (int64, error)
}

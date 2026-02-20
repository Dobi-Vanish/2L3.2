package repository

import (
	"context"
	"errors"
	"fmt"
	"shortener/internal/logger"
	"shortener/internal/model"

	"github.com/wb-go/wbf/dbpg/pgx-driver"
)

type PostgresRepository struct {
	db *pgxdriver.Postgres
}

func NewPostgresRepository(db *pgxdriver.Postgres) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveLink(ctx context.Context, link *model.Link) error {
	query := `INSERT INTO links (short_url, original_url, custom_alias, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, link.ShortURL, link.OriginalURL, link.CustomAlias, link.CreatedAt)
	if err != nil {
		logger.Error("failed to save link", "error", err)
		return fmt.Errorf("save link: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetLink(ctx context.Context, shortURL string) (*model.Link, error) {
	query := `SELECT short_url, original_url, custom_alias, created_at FROM links WHERE short_url = $1`
	row := r.db.QueryRow(ctx, query, shortURL)
	var link model.Link
	err := row.Scan(&link.ShortURL, &link.OriginalURL, &link.CustomAlias, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, err) {
			return nil, nil
		}
		logger.Error("failed to get link", "error", err)
		return nil, fmt.Errorf("get link: %w", err)
	}
	return &link, nil
}

func (r *PostgresRepository) LinkExists(ctx context.Context, shortURL string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM links WHERE short_url = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, shortURL).Scan(&exists)
	if err != nil {
		logger.Error("failed to check link existence", "error", err)
		return false, fmt.Errorf("link exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) SaveAnalytics(ctx context.Context, a *model.Analytics) error {
	query := `INSERT INTO analytics (short_url, timestamp, user_agent, referer) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, a.ShortURL, a.Timestamp, a.UserAgent, a.Referer)
	if err != nil {
		logger.Error("failed to save analytics", "error", err)
		return fmt.Errorf("save analytics: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetAnalytics(ctx context.Context, shortURL string, limit int) ([]model.Analytics, error) {
	query := `SELECT id, short_url, timestamp, user_agent, referer FROM analytics WHERE short_url = $1 ORDER BY timestamp DESC LIMIT $2`
	rows, err := r.db.Query(ctx, query, shortURL, limit)
	if err != nil {
		logger.Error("failed to get analytics", "error", err)
		return nil, fmt.Errorf("get analytics: %w", err)
	}
	defer rows.Close()

	var analytics []model.Analytics
	for rows.Next() {
		var a model.Analytics
		if err := rows.Scan(&a.ID, &a.ShortURL, &a.Timestamp, &a.UserAgent, &a.Referer); err != nil {
			logger.Error("failed to scan analytics row", "error", err)
			return nil, fmt.Errorf("scan analytics: %w", err)
		}
		analytics = append(analytics, a)
	}
	if err := rows.Err(); err != nil {
		logger.Error("rows iteration error", "error", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return analytics, nil
}

func (r *PostgresRepository) CountAnalytics(ctx context.Context, shortURL string) (int64, error) {
	query := `SELECT COUNT(*) FROM analytics WHERE short_url = $1`
	var count int64
	err := r.db.QueryRow(ctx, query, shortURL).Scan(&count)
	if err != nil {
		logger.Error("failed to count analytics", "error", err)
		return 0, fmt.Errorf("count analytics: %w", err)
	}
	return count, nil
}

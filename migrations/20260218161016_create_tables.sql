-- +goose Up
CREATE TABLE IF NOT EXISTS links (
    short_url VARCHAR(50) PRIMARY KEY,
    original_url TEXT NOT NULL,
    custom_alias VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS analytics (
    id BIGSERIAL PRIMARY KEY,
    short_url VARCHAR(50) NOT NULL REFERENCES links(short_url) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_agent TEXT,
    referer TEXT
    );

CREATE INDEX idx_analytics_short_url ON analytics(short_url);
CREATE INDEX idx_analytics_timestamp ON analytics(timestamp);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS analytics;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

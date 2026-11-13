-- +goose Up
-- email is already indexed by the UNIQUE constraint; a second index only slows writes.
DROP INDEX IF EXISTS idx_users_email;

CREATE INDEX idx_users_created_at_id ON users (created_at DESC, id DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_users_created_at_id;

CREATE INDEX idx_users_email ON users (email);

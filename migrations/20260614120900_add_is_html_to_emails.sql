-- +goose Up
ALTER TABLE emails ADD COLUMN is_html BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE emails DROP COLUMN is_html;

-- +goose Up
ALTER TABLE html_patterns RENAME TO html_templates;

-- +goose Down
ALTER TABLE html_templates RENAME TO html_patterns;

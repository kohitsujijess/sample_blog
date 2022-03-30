-- +goose Up
-- +goose StatementBegin
ALTER TABLE articles DROP COLUMN description;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE articles ADD COLUMN description text AFTER title;
-- +goose StatementEnd

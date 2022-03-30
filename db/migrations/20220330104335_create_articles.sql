-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles (
    id int NOT NULL,
    title text,
    description text,
    body text,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE articles;
-- +goose StatementEnd

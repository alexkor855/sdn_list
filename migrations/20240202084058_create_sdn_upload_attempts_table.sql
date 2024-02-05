-- +goose Up
-- +goose StatementBegin
CREATE TABLE sdn_upload_attempts (
    id SERIAL PRIMARY KEY,
    status SMALLINT NOT NULL DEFAULT(1),
    publish_date TIMESTAMP,
    started_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sdn_upload_attempts;
-- +goose StatementEnd

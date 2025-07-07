-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS examples(
  example_id serial primary key,
  example_name VARCHAR not null unique,
  example_description VARCHAR,
  example_created_at timestamp default (now() at time zone 'utc'),
  example_updated_at timestamp default (now() at time zone 'utc')
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS examples;
-- +goose StatementEnd

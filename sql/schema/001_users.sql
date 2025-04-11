-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY,
  created_at timestamp not null,
  updated_at timestamp not null,
  email text not null UNIQUE
);
-- +goose Down
DROP TABLE users;

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pickup_points
(
    id      bigserial primary key not null,
    name    text                  not null,
    address text                  not null,
    contact text                  not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pickup_points;
-- +goose StatementEnd

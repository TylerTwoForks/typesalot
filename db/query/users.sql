-- name: CreateUser :exec
INSERT INTO users (name) VALUES (?);

-- name: ListUsers :many
SELECT * FROM users;
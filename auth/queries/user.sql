-- name: GetUserByLogin :one
SELECT id, login, password_hash, role
FROM "user"
WHERE login = $1;

-- name: CreateUser :one
INSERT INTO "user" (login, password_hash, role)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateUser :exec
UPDATE "user"
SET login = $2, password_hash = $3, role = $4
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "user" WHERE id = $1;

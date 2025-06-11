-- name: GetUserByForeignID :one
SELECT *
FROM users
WHERE foreign_id = $1
LIMIT 1;
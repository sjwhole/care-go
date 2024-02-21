-- name: GetParent :one
SELECT * FROM parents WHERE user_id = ? LIMIT 1;

-- name: ListParents :many
SELECT * FROM parents WHERE user_id = ? ORDER BY name;
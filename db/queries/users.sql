-- name: GetUserById :one
SELECT * FROM users WHERE id = ? LIMIT 1;


-- name: GetUserByKakaoId :one
SELECT * FROM users WHERE kakao_id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :execresult
INSERT INTO users (
  kakao_id, name, phone_no
) VALUES (
  ?, ?, ?
);


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;


-- name: UpdateUser :exec
UPDATE users
SET name = ?, phone_no = ?
WHERE id = ?;


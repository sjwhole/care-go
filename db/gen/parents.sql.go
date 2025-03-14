// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: parents.sql

package db

import (
	"context"
	"database/sql"
)

const getParent = `-- name: GetParent :one
SELECT id, created_at, updated_at, deleted_at, user_id, name, phone_no FROM parents WHERE user_id = ? LIMIT 1
`

func (q *Queries) GetParent(ctx context.Context, userID sql.NullInt64) (Parent, error) {
	row := q.db.QueryRowContext(ctx, getParent, userID)
	var i Parent
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.Name,
		&i.PhoneNo,
	)
	return i, err
}

const listParents = `-- name: ListParents :many
SELECT id, created_at, updated_at, deleted_at, user_id, name, phone_no FROM parents WHERE user_id = ? ORDER BY name
`

func (q *Queries) ListParents(ctx context.Context, userID sql.NullInt64) ([]Parent, error) {
	rows, err := q.db.QueryContext(ctx, listParents, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Parent
	for rows.Next() {
		var i Parent
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.UserID,
			&i.Name,
			&i.PhoneNo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

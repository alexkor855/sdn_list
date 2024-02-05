// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package upload_attempts_queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getLastSdnUploadAttempt = `-- name: GetLastSdnUploadAttempt :one
SELECT id, status, publish_date, started_at 
FROM sdn_upload_attempts 
ORDER BY id DESC
LIMIT 1
`

func (q *Queries) GetLastSdnUploadAttempt(ctx context.Context) (SdnUploadAttempt, error) {
	row := q.db.QueryRow(ctx, getLastSdnUploadAttempt)
	var i SdnUploadAttempt
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.PublishDate,
		&i.StartedAt,
	)
	return i, err
}

const getSdnUploadAttempt = `-- name: GetSdnUploadAttempt :one
SELECT id, status, publish_date, started_at FROM sdn_upload_attempts WHERE id = $1
`

func (q *Queries) GetSdnUploadAttempt(ctx context.Context, id int32) (SdnUploadAttempt, error) {
	row := q.db.QueryRow(ctx, getSdnUploadAttempt, id)
	var i SdnUploadAttempt
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.PublishDate,
		&i.StartedAt,
	)
	return i, err
}

const insertSdnUploadAttempt = `-- name: InsertSdnUploadAttempt :one
INSERT INTO sdn_upload_attempts (status) VALUES (1) RETURNING id
`

func (q *Queries) InsertSdnUploadAttempt(ctx context.Context) (int32, error) {
	row := q.db.QueryRow(ctx, insertSdnUploadAttempt)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const updateSdnUploadAttempt = `-- name: UpdateSdnUploadAttempt :exec
UPDATE sdn_upload_attempts
SET status = 2, publish_date = $2
WHERE id = $1
`

type UpdateSdnUploadAttemptParams struct {
	ID          int32            `json:"id"`
	PublishDate pgtype.Timestamp `json:"publish_date"`
}

func (q *Queries) UpdateSdnUploadAttempt(ctx context.Context, arg UpdateSdnUploadAttemptParams) error {
	_, err := q.db.Exec(ctx, updateSdnUploadAttempt, arg.ID, arg.PublishDate)
	return err
}
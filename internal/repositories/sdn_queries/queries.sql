-- name: InsertSdn :one
INSERT INTO sdn_list (uid, first_name, last_name) VALUES ($1, $2, $3) RETURNING id;

-- name: GetSdnByUid :many
SELECT * FROM sdn_list WHERE uid = $1;

-- name: GetSdnById :one
SELECT * FROM sdn_list WHERE id = $1;

-- name: GetSdnByUidAndName :one
SELECT * FROM sdn_list WHERE uid = $1 AND first_name = $2 AND last_name = $3;

-- name: DeleteOrder :exec
DELETE FROM sdn_list WHERE id = $1;

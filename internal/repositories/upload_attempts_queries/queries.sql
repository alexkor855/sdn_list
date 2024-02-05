-- name: InsertSdnUploadAttempt :one
INSERT INTO sdn_upload_attempts (status) VALUES (1) RETURNING id;

-- name: UpdateSdnUploadAttempt :exec
UPDATE sdn_upload_attempts
SET status = 2, publish_date = $2
WHERE id = $1;

-- name: GetSdnUploadAttempt :one
SELECT * FROM sdn_upload_attempts WHERE id = $1;

-- name: GetLastSdnUploadAttempt :one
SELECT * 
FROM sdn_upload_attempts 
ORDER BY id DESC
LIMIT 1;

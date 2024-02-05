// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package upload_attempts_queries

import (
	"context"
)

type Querier interface {
	GetLastSdnUploadAttempt(ctx context.Context) (SdnUploadAttempt, error)
	GetSdnUploadAttempt(ctx context.Context, id int32) (SdnUploadAttempt, error)
	InsertSdnUploadAttempt(ctx context.Context) (int32, error)
	UpdateSdnUploadAttempt(ctx context.Context, arg UpdateSdnUploadAttemptParams) error
}

var _ Querier = (*Queries)(nil)

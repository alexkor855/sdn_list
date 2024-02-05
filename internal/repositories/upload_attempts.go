package repositories

import (
	"context"
	"errors"
	"sdn_list/internal/entities"
	"sdn_list/internal/repositories/upload_attempts_queries"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UploadAttemptsRepository struct {
	dbPool *pgxpool.Pool
}

func NewUploadAttemptsRepository(dbPool *pgxpool.Pool) *UploadAttemptsRepository {
	return &UploadAttemptsRepository{dbPool: dbPool}
}

func (r *UploadAttemptsRepository) GetLastSdnUploadAttempt(ctx context.Context) (*entities.UploadAttempt, error) {
	queries := upload_attempts_queries.New(r.dbPool)

	uploadAttempt, err := queries.GetLastSdnUploadAttempt(ctx)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return convertFromRowToEntity(uploadAttempt), nil
	}
}

func convertFromRowToEntity(raw upload_attempts_queries.SdnUploadAttempt) *entities.UploadAttempt {
	return &entities.UploadAttempt{
		Id:          raw.ID,
		Status:      entities.UploadStatus(raw.Status),
		PublishDate: raw.PublishDate.Time,
		StartedAt:   raw.StartedAt.Time,
	}
}

func (r *UploadAttemptsRepository) Create(ctx context.Context) (int32, error) {
	queries := upload_attempts_queries.New(r.dbPool)

	attemptId, err := queries.InsertSdnUploadAttempt(ctx)
	if err != nil {
		return 0, err
	}

	return attemptId, nil
}

func (r *UploadAttemptsRepository) UpdateSuccessAttempt(ctx context.Context, id int32, publishDate string) error {
	queries := upload_attempts_queries.New(r.dbPool)

	publishDateTime, _ := time.Parse("02/01/2006", publishDate)

	pgPublishDate := pgtype.Timestamp{}
	pgPublishDate.Scan(publishDateTime)

	return queries.UpdateSdnUploadAttempt(ctx, upload_attempts_queries.UpdateSdnUploadAttemptParams{
		ID:	id,
		PublishDate: pgPublishDate,
	})
}

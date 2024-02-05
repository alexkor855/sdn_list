package services

import (
	"context"
	"sdn_list/internal/entities"
)

type UploadAttemptsRepository interface {
	Create(ctx context.Context) (int32, error)
	UpdateSuccessAttempt(ctx context.Context, id int32, publishDate string) error
	GetLastSdnUploadAttempt(ctx context.Context) (*entities.UploadAttempt, error)
}

type SdnRepository interface {
	SaveAll(ctx context.Context, input <-chan Person) error
}

type SdnProvider interface {
	GetAll(ctx context.Context) (<-chan Person, <-chan MetaData, error)
}

type UploadService struct {
	name                     string
	uploadAttemptsRepository UploadAttemptsRepository
	sdnRepository            SdnRepository
	sdnProvider              SdnProvider
}

func NewUploadService(uploadAttemptsRepository UploadAttemptsRepository, sdnRepository SdnRepository, sdnProvider SdnProvider) *UploadService {
	return &UploadService{
		name:                     "upload service",
		uploadAttemptsRepository: uploadAttemptsRepository,
		sdnRepository:            sdnRepository,
		sdnProvider:              sdnProvider,
	}
}

func (us *UploadService) IsLastUploadSuccessful(ctx context.Context) (bool, error) {

	lastUpload, err := us.uploadAttemptsRepository.GetLastSdnUploadAttempt(ctx)
	if err != nil {
		return false, err
	}

	if lastUpload == nil {
		return false, nil
	}

	return lastUpload.Status == entities.Successful, nil
}

func (us *UploadService) Upload(ctx context.Context) error {

	// получить историю загрузок
	// сравнить даты публикации, если не изменилась, то ничего не делать
	// ....

	attemptId, err := us.uploadAttemptsRepository.Create(ctx)
	if err != nil {
		return err
	}

	input, metaData, err := us.sdnProvider.GetAll(ctx)
	if err != nil {
		return err
	}

	err = us.sdnRepository.SaveAll(ctx, input)
	if err != nil {
		return err
	}

	metaInfo := <-metaData
	us.uploadAttemptsRepository.UpdateSuccessAttempt(ctx, attemptId, metaInfo.PublishDate)

	return nil
}

type Person struct {
	Uid       int    
	FirstName string
	LastName  string
}

type MetaData struct {
	PublishDate string
}

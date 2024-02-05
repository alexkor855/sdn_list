package entities

import "time"

type UploadAttempt struct {
	Id          int32
	Status      UploadStatus
	PublishDate time.Time
	StartedAt   time.Time
}

type UploadStatus int16

const (
	Failed UploadStatus = iota + 1
	Successful
)
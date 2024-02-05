package services

import (
	"context"
	"sdn_list/internal/entities"
)

type SdnSearchRepository interface {
	SearchStrong(ctx context.Context, name string) ([]entities.Person, error)
	SearchWeak(ctx context.Context, name string) ([]entities.Person, error)
}

type SearchService struct {
	name          string
	sdnRepository SdnSearchRepository
}

func NewSearchService(sdnRepository SdnSearchRepository) *SearchService {
	return &SearchService{
		name:          "search service",
		sdnRepository: sdnRepository,
	}
}

func (s *SearchService) Search(ctx context.Context, name string, isStrongParam bool) ([]entities.Person, error) {
	if isStrongParam {
		return s.searchStrong(ctx, name)
	} else {
		return s.searchWeak(ctx, name)
	}
}

func (s *SearchService) searchStrong(ctx context.Context, name string) ([]entities.Person, error) {
	return s.sdnRepository.SearchStrong(ctx, name)
}

func (s *SearchService) searchWeak(ctx context.Context, name string) ([]entities.Person, error) {
	return s.sdnRepository.SearchWeak(ctx, name)
}

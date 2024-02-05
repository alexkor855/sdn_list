package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sdn_list/internal/entities"
	"strings"
)

type SearchService interface {
	Search(ctx context.Context, name string, isStrongParam bool) ([]entities.Person, error)
}

type SearchHandler struct {
	searchService SearchService
}

func NewSearchHandler(searchService SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

type ItemResponse struct {
	Uid       int    `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *SearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isStrongParam := false
	nameParam := ""

	params := r.URL.Query()
	for param, val := range params {
		if strings.ToLower(param) == "type" {
			isStrongParam = strings.ToLower(val[0]) == "strong"
		}
		if strings.ToLower(param) == "name" {
			nameParam = val[0]
		}
	}

	persons, err := h.searchService.Search(ctx, nameParam, isStrongParam)
	if err != nil {
		GetErrorResponse(w, "search", err, http.StatusInternalServerError)
	}

	searchResponse := convertFromEntityToResponse(persons)
	raw, err := json.Marshal(searchResponse)
	if err != nil {
		GetErrorResponse(w, "search", err, http.StatusInternalServerError)
	}

	GetSuccessResponseWithBody(w, raw)
}

func convertFromEntityToResponse(persons []entities.Person) []ItemResponse {
	result := make([]ItemResponse, len(persons))
	for i, person := range persons {
		result[i] = ItemResponse{
			Uid:       person.Uid,
			FirstName: person.FirstName,
			LastName:  person.LastName,
		}
	}
	return result
}

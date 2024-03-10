package repository

import "go-image-api/internal/entity"

type HistoryRepository struct {
	Repository[entity.History]
}

func NewHistoryRepository() *HistoryRepository {
	return new(HistoryRepository)
}

package services

import (
	"tranee_service/internal/logging"
	"tranee_service/models"
	"tranee_service/repositories"
)

type HobbyService struct {
	repository *repositories.Repository
	logger     logging.Logger
}

func NewHobbyService(repository *repositories.Repository, logger logging.Logger) *HobbyService {
	return &HobbyService{repository: repository, logger: logger}
}

func (h *HobbyService) CreateHobby(hobby *models.Hobby) (int, error) {
	return h.repository.AppHobbies.CreateHobby(hobby)
}

func (h *HobbyService) GetHobbies() ([]models.ResponseHobby, error) {
	return h.repository.AppHobbies.GetHobbies()
}

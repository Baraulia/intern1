package services

import (
	"tranee_service/internal/logging"
	"tranee_service/models"
	"tranee_service/repositories"
)

type UserService struct {
	repository *repositories.Repository
	logger     logging.Logger
}

func NewUserService(repository *repositories.Repository, logger logging.Logger) *UserService {
	return &UserService{repository: repository, logger: logger}
}

func (u *UserService) CreateUser(user *models.User) (int, error) {
	return u.repository.AppUsers.CreateUser(user)
}

func (u *UserService) GetUserById(userId int) (*models.ResponseUser, error) {
	return u.repository.AppUsers.GetUserById(userId)
}

func (u *UserService) GetUsers(options *models.Options) ([]models.ResponseUser, int, error) {
	return u.repository.AppUsers.GetUsers(options)
}

func (u *UserService) ChangeUser(user *models.User, userId int) error {
	return u.repository.AppUsers.ChangeUser(user, userId)
}

func (u *UserService) DeleteUser(userId int) error {
	return u.repository.AppUsers.DeleteUser(userId)
}
func (u *UserService) GetHobbyByUserId(userId int) ([]int, error) {
	return u.repository.AppHobbies.GetHobbyByUserId(userId)
}

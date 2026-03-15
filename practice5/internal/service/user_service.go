package service

import (
	"practice5/internal/model"
	"practice5/internal/repository"
)

type UserService interface {
	GetPaginatedUsers(page, pageSize int, filters map[string]string, orderBy string) (model.PaginatedResponse, error)
	GetUserByID(id int) (*model.User, error)
	CreateUser(req model.CreateUserRequest) (*model.User, error)
	UpdateUser(id int, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(id int) error
	GetCommonFriends(user1ID, user2ID int) ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetPaginatedUsers(page, pageSize int, filters map[string]string, orderBy string) (model.PaginatedResponse, error) {
	return s.repo.GetPaginatedUsers(page, pageSize, filters, orderBy)
}

func (s *userService) GetUserByID(id int) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) CreateUser(req model.CreateUserRequest) (*model.User, error) {
	return s.repo.CreateUser(req)
}

func (s *userService) UpdateUser(id int, req model.UpdateUserRequest) (*model.User, error) {
	return s.repo.UpdateUser(id, req)
}

func (s *userService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) GetCommonFriends(user1ID, user2ID int) ([]model.User, error) {
	return s.repo.GetCommonFriends(user1ID, user2ID)
}

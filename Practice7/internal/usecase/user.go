package usecase

import (
	"fmt"
	"practice7/internal/entity"
	"practice7/internal/usecase/repo"
	"practice7/utils"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	user, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}

	sessionID := uuid.New().String()
	return user, sessionID, nil
}

func (u *UserUseCase) LoginUser(user *entity.LoginUserDTO) (string, error) {
	userFromRepo, err := u.repo.LoginUser(user)
	if err != nil {
		return "", fmt.Errorf("user from repo: %w", err)
	}

	if !utils.CheckPassword(userFromRepo.Password, user.Password) {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("generate JWT: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) GetMe(userID uuid.UUID) (*entity.User, error) {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("get me: %w", err)
	}

	return user, nil
}

func (u *UserUseCase) PromoteUser(targetUserID uuid.UUID) (*entity.User, error) {
	user, err := u.repo.GetByID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("find target user: %w", err)
	}

	user.Role = "admin"

	updatedUser, err := u.repo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("promote user: %w", err)
	}

	return updatedUser, nil
}

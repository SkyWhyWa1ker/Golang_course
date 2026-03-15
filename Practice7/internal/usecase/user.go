package usecase

import (
	"fmt"
	"net/http"
	"practice7/internal/entity"
	"practice7/internal/usecase"
	"practice7/internal/usecase/repo"
	"practice7/utils"
	"time"

	"github.com/gin-gonic/gin"
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

func (r *UserRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO

	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}

	role := createUserDTO.Role
	if role == "" {
		role = "user"
	}

	user := entity.User{
		ID:       uuid.New(),
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     role,
		Verified: false,
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "user registered successfully",
		"session_id": sessionID,
		"user": gin.H{
			"id":       createdUser.ID,
			"username": createdUser.Username,
			"email":    createdUser.Email,
			"role":     createdUser.Role,
			"verified": createdUser.Verified,
		},
	})
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

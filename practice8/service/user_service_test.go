package service

import (
	"errors"
	"practice8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{
		ID:    1,
		Name:  "Emir",
		Email: "emir@example.com",
	}

	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{
		ID:    2,
		Name:  "New User",
		Email: "new@example.com",
	}

	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)

	assert.NoError(t, err)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	existingUser := &repository.User{
		ID:    1,
		Name:  "Existing",
		Email: "existing@example.com",
	}

	newUser := &repository.User{
		ID:    2,
		Name:  "New",
		Email: "existing@example.com",
	}

	mockRepo.EXPECT().GetByEmail("existing@example.com").Return(existingUser, nil)

	err := userService.RegisterUser(newUser, "existing@example.com")

	assert.Error(t, err)
	assert.EqualError(t, err, "user with this email already exists")
}

func TestRegisterUser_NewUserSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &repository.User{
		ID:    2,
		Name:  "New",
		Email: "new@example.com",
	}

	mockRepo.EXPECT().GetByEmail("new@example.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(nil)

	err := userService.RegisterUser(newUser, "new@example.com")

	assert.NoError(t, err)
}

func TestRegisterUser_CreateUserRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &repository.User{
		ID:    2,
		Name:  "New",
		Email: "new@example.com",
	}

	mockRepo.EXPECT().GetByEmail("new@example.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(errors.New("create failed"))

	err := userService.RegisterUser(newUser, "new@example.com")

	assert.Error(t, err)
	assert.EqualError(t, err, "create failed")
}

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	err := userService.UpdateUserName(1, "")

	assert.Error(t, err)
	assert.EqualError(t, err, "name cannot be empty")
}

func TestUpdateUserName_UserNotFoundOrRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(99).Return(nil, errors.New("user not found"))

	err := userService.UpdateUserName(99, "NewName")

	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{
		ID:    5,
		Name:  "OldName",
		Email: "user@example.com",
	}

	mockRepo.EXPECT().GetUserByID(5).Return(user, nil)
	mockRepo.EXPECT().
		UpdateUser(gomock.Any()).
		DoAndReturn(func(updatedUser *repository.User) error {
			assert.Equal(t, "NewName", updatedUser.Name)
			return nil
		})

	err := userService.UpdateUserName(5, "NewName")

	assert.NoError(t, err)
	assert.Equal(t, "NewName", user.Name)
}

func TestUpdateUserName_UpdateUserFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{
		ID:    7,
		Name:  "OldName",
		Email: "user7@example.com",
	}

	mockRepo.EXPECT().GetUserByID(7).Return(user, nil)
	mockRepo.EXPECT().
		UpdateUser(gomock.Any()).
		DoAndReturn(func(updatedUser *repository.User) error {
			assert.Equal(t, "UpdatedName", updatedUser.Name)
			return errors.New("update failed")
		})

	err := userService.UpdateUserName(7, "UpdatedName")

	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
	assert.Equal(t, "UpdatedName", user.Name)
}

func TestDeleteUser_AttemptToDeleteAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	err := userService.DeleteUser(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "it is not allowed to delete admin user")
}

func TestDeleteUser_SuccessfulDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	deleted := false

	mockRepo.EXPECT().
		DeleteUser(2).
		DoAndReturn(func(id int) error {
			deleted = true
			return nil
		})

	err := userService.DeleteUser(2)

	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(3).Return(errors.New("delete failed"))

	err := userService.DeleteUser(3)

	assert.Error(t, err)
	assert.EqualError(t, err, "delete failed")
}

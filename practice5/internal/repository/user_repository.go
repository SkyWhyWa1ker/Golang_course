package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"practice5/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetPaginatedUsers(page, pageSize int, filters map[string]string, orderBy string) (model.PaginatedResponse, error)
	GetUserByID(id int) (*model.User, error)
	CreateUser(req model.CreateUserRequest) (*model.User, error)
	UpdateUser(id int, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(id int) error
	GetCommonFriends(user1ID, user2ID int) ([]model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetPaginatedUsers(page, pageSize int, filters map[string]string, orderBy string) (model.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}

	offset := (page - 1) * pageSize

	allowedOrderBy := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"gender":     "gender",
		"birth_date": "birth_date",
	}

	if orderBy == "" {
		orderBy = "id"
	}
	orderColumn, ok := allowedOrderBy[orderBy]
	if !ok {
		orderColumn = "id"
	}

	baseQuery := " FROM users WHERE 1=1"
	args := []interface{}{}
	argPos := 1

	if value, ok := filters["id"]; ok && value != "" {
		id, err := strconv.Atoi(value)
		if err != nil {
			return model.PaginatedResponse{}, errors.New("invalid id filter")
		}
		baseQuery += fmt.Sprintf(" AND id = $%d", argPos)
		args = append(args, id)
		argPos++
	}

	if value, ok := filters["name"]; ok && value != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+value+"%")
		argPos++
	}

	if value, ok := filters["email"]; ok && value != "" {
		baseQuery += fmt.Sprintf(" AND email ILIKE $%d", argPos)
		args = append(args, "%"+value+"%")
		argPos++
	}

	if value, ok := filters["gender"]; ok && value != "" {
		baseQuery += fmt.Sprintf(" AND gender = $%d", argPos)
		args = append(args, strings.ToLower(value))
		argPos++
	}

	if value, ok := filters["birth_date"]; ok && value != "" {
		_, err := time.Parse("2006-01-02", value)
		if err != nil {
			return model.PaginatedResponse{}, errors.New("invalid birth_date filter, use YYYY-MM-DD")
		}
		baseQuery += fmt.Sprintf(" AND birth_date = $%d", argPos)
		args = append(args, value)
		argPos++
	}

	countQuery := "SELECT COUNT(*)" + baseQuery
	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return model.PaginatedResponse{}, err
	}

	dataQuery := `
		SELECT id, name, email, gender, birth_date
	` + baseQuery + fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", orderColumn, argPos, argPos+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return model.PaginatedResponse{}, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return model.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	return model.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	query := `
		SELECT id, name, email, gender, birth_date
		FROM users
		WHERE id = $1
	`

	var u model.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) CreateUser(req model.CreateUserRequest) (*model.User, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, errors.New("invalid birth_date format, use YYYY-MM-DD")
	}

	query := `
		INSERT INTO users (name, email, gender, birth_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, gender, birth_date
	`

	var u model.User
	err = r.db.QueryRow(query, req.Name, req.Email, strings.ToLower(req.Gender), birthDate).
		Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) UpdateUser(id int, req model.UpdateUserRequest) (*model.User, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, errors.New("invalid birth_date format, use YYYY-MM-DD")
	}

	query := `
		UPDATE users
		SET name = $1, email = $2, gender = $3, birth_date = $4
		WHERE id = $5
		RETURNING id, name, email, gender, birth_date
	`

	var u model.User
	err = r.db.QueryRow(query, req.Name, req.Email, strings.ToLower(req.Gender), birthDate, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) DeleteUser(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) GetCommonFriends(user1ID, user2ID int) ([]model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON u.id = uf1.friend_id
		JOIN user_friends uf2 ON u.id = uf2.friend_id
		WHERE uf1.user_id = $1
		  AND uf2.user_id = $2
		ORDER BY u.id
	`

	rows, err := r.db.Query(query, user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

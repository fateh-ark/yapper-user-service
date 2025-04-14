package repositories

import (
	"context"

	"fateh-ark/yapper-user-service/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
}

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Inserts a new user into the database.
func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (username, email, display_name, profile_image_url)
		 VALUES ($1, $2, $3, $4)`,
		user.Username, user.Email, user.DisplayName, user.ProfileImageURL,
	)
	return err
}

// Retrieves a user by their ID.
func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, display_name, profile_image_url, created_at, updated_at
		 FROM users
		 WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.DisplayName, &user.ProfileImageURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Retrieves a user by their username.
func (r *userRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, display_name, profile_image_url, created_at, updated_at
		 FROM users
		 WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.DisplayName, &user.ProfileImageURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Retrieves a user by their email.
func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, display_name, profile_image_url, created_at, updated_at
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.DisplayName, &user.ProfileImageURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Updates an existing user in the database.
func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users
		 SET username = $2, email = $3, display_name = $4, profile_image_url = $5
		 WHERE id = $1`,
		user.ID, user.Username, user.Email, user.DisplayName, user.ProfileImageURL,
	)
	return err
}

// Deletes a user from the database by their ID.
func (r *userRepositoryImpl) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM users
		 WHERE id = $1`,
		id,
	)
	return err
}

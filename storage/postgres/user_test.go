package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	pb "Auth-Service/genproto/users"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestRegister(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	ctx := context.Background()
	req := &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		FullName: "Test User",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(req.Username, req.Email, req.Password, req.FullName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("12345", time.Now()))

	resp, err := repo.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "Test User", resp.FullName)
	assert.Equal(t, "12345", resp.Id)
}

func TestLogin(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	ctx := context.Background()
	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "password",
	}

	mock.ExpectQuery("SELECT id, username, email, full_name, created_at FROM users WHERE username = \\$1 AND password = \\$2").
		WithArgs(req.Username, req.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "full_name", "created_at"}).AddRow("12345", "testuser", "test@example.com", "Test User", time.Now()))

	resp, err := repo.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "Test User", resp.FullName)
	assert.Equal(t, "12345", resp.Id)
}

func TestGetUserByID(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	ctx := context.Background()
	userID := "12345"

	mock.ExpectQuery("SELECT username, email, password, full_name, bio, countries_visited FROM users WHERE id = \\$1 AND deleted_at IS NULL").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password", "full_name", "bio", "countries_visited"}).
			AddRow("testuser", "test@example.com", "password", "Test User", "Bio", 5))

	resp, err := repo.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "password", resp.Password)
	assert.Equal(t, "Test User", resp.FullName)
	assert.Equal(t, "Bio", resp.Bio)
	assert.Equal(t, int32(5), resp.CountriesVisited)
}

func TestProfile(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	ctx := context.Background()
	req := &pb.ProfileRequest{
		UserId: "12345",
	}

	mock.ExpectQuery("SELECT id, username, email, full_name, bio, countries_visited, created_at, updated_at FROM users WHERE id=\\$1 AND deleted_at IS NULL").
		WithArgs(req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "full_name", "bio", "countries_visited", "created_at", "updated_at"}).
			AddRow("12345", "testuser", "test@example.com", "Test User", "Bio", 5, time.Now(), time.Now()))

	resp, err := repo.Profile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "12345", resp.Id)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "Test User", resp.FullName)
	assert.Equal(t, "Bio", resp.Bio)
	assert.Equal(t, int32(5), resp.CountriesVisited)
}

func TestUpdateProfile(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	ctx := context.Background()
	req := &pb.UpdateProfileRequest{
		Id:               "12345",
		FullName:         "Updated User",
		Bio:              "Updated Bio",
		CountriesVisited: 10,
	}

	mock.ExpectQuery("UPDATE users SET full_name = \\$1, bio = \\$2, countries_visited = \\$3, updated_at = \\$4 WHERE id = \\$5 RETURNING id, username, email, full_name, bio, countries_visited, updated_at").
		WithArgs(req.FullName, req.Bio, req.CountriesVisited, sqlmock.AnyArg(), req.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "full_name", "bio", "countries_visited", "updated_at"}).
			AddRow("12345", "testuser", "test@example.com", "Updated User", "Updated Bio", 10, time.Now()))

	resp, err := repo.UpdateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "12345", resp.Id)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "Updated User", resp.FullName)
	assert.Equal(t, "Updated Bio", resp.Bio)
	assert.Equal(t, int32(10), resp.CountriesVisited)
}

func TestDeleteUser(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	ctx := context.Background()
	req := &pb.DeleteUserRequest{
		Id: "12345",
	}

	mock.ExpectExec("UPDATE users SET deleted_at=CURRENT_TIMESTAMP WHERE id=\\$1 AND deleted_at IS NULL").
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := repo.DeleteUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.StatusUser)
}

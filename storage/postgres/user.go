package postgres

import (
	pb "Auth-Service/genproto/users"
	storage "Auth-Service/help"
	"context"
	"database/sql"
	"fmt"
	"net/smtp"
	"time"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (repo *UserRepository) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if repo.Db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var id, createdAt string
	err := repo.Db.QueryRowContext(ctx,
		`INSERT INTO users (username, email, password, full_name) 
         VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		request.Username, request.Email, request.Password, request.FullName,
	).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	response := &pb.RegisterResponse{
		Id:        id,
		Username:  request.Username,
		Email:     request.Email,
		FullName:  request.FullName,
		CreatedAt: createdAt,
	}

	return response, nil
}

func (repo *UserRepository) Login(ctx context.Context, request *pb.LoginRequest) (*pb.RegisterResponse, error) {
	var loginUser pb.RegisterResponse
	err := repo.Db.QueryRowContext(ctx, 
		"SELECT id, username, email, full_name, created_at FROM users WHERE username = $1 AND password = $2", 
		request.Username, request.Password,
	).Scan(&loginUser.Id, &loginUser.Username, &loginUser.Email, &loginUser.FullName, &loginUser.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &loginUser, nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, id string) (*pb.UserInfo, error) {
	user := &pb.UserInfo{Id: id}

	query := `
	SELECT
		username,
		email,
		password,
		full_name,
		bio,
		countries_visited
	FROM
		users
	WHERE
		id = $1 AND deleted_at IS NULL
	`
	row := repo.Db.QueryRowContext(ctx, query, id)

	var bio sql.NullString
	err := row.Scan(&user.Username, &user.Email, &user.Password, &user.FullName, &bio, &user.CountriesVisited)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if bio.Valid {
		user.Bio = bio.String
	} else {
		user.Bio = ""
	}

	return user, nil
}

func (repo *UserRepository) Profile(ctx context.Context, request *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	var user pb.ProfileResponse
	err := repo.Db.QueryRowContext(
		ctx,
		"SELECT id, username, email, full_name, bio, countries_visited, created_at, updated_at FROM users WHERE id=$1 AND deleted_at IS NULL",
		request.UserId,
	).Scan(&user.Id, &user.Username, &user.Email, &user.FullName, &user.Bio, &user.CountriesVisited, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateProfile(ctx context.Context, request *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	query := `UPDATE users 
			  SET full_name = $1, bio = $2, countries_visited = $3, updated_at = $4 
			  WHERE id = $5
			  RETURNING id, username, email, full_name, bio, countries_visited, updated_at`

	row := repo.Db.QueryRowContext(ctx, query,
		request.FullName, request.Bio, request.CountriesVisited, time.Now(), request.Id)

	response := &pb.UpdateProfileResponse{}

	err := row.Scan(
		&response.Id,
		&response.Username,
		&response.Email,
		&response.FullName,
		&response.Bio,
		&response.CountriesVisited,
		&response.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating profile: %v", err)
	}

	return response, nil
}

func (repo *UserRepository) GetUsers(ctx context.Context, request *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	var (
		params = make(map[string]interface{})
		arr    []interface{}
	)
	filter := ""
	if request.Limit > 0 {
		params["limit"] = request.Limit
		filter += " LIMIT :limit "
	}
	if request.Offset > 0 {
		params["offset"] = request.Offset
		filter += " OFFSET :offset "
	}

	query := "SELECT id, username, full_name, countries_visited FROM users WHERE deleted_at IS NULL"
	query = query + filter
	query, arr = storage.ReplaceQueryParams(query, params)
	rows, err := repo.Db.QueryContext(ctx, query, arr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.Users
	for rows.Next() {
		var user pb.Users
		err := rows.Scan(&user.Id, &user.Username, &user.FullName, &user.CountriesVisited)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &pb.GetUsersResponse{Users: users, Limit: request.Limit, Total: int32(len(users))}, nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, err := repo.Db.ExecContext(
		ctx,
		"UPDATE users SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL",
		request.Id,
	)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{StatusUser: true}, nil
}

func (repo *UserRepository) ResetPassword(ctx context.Context, request *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	emailBody := "Click the link to reset your password: https://your-domain.com/reset-password"

	err := SendEmail(request.Email, "Password Reset Instructions", emailBody)
	if err != nil {
		return nil, fmt.Errorf("error sending reset email: %v", err)
	}

	return &pb.ResetPasswordResponse{
		Message: "Password reset instructions sent to your email",
	}, nil
}

func SendEmail(to, subject, body string) error {
	from := "your-email@example.com"
	password := "your-email-password"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	fmt.Println("Email sent to:", to)
	return nil
}

func (repo *UserRepository) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	query := `UPDATE users
			SET token = NULL
			WHERE id = $1`

	result, err := repo.Db.ExecContext(ctx, query, request.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to update user token: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found or token update failed")
	}

	return &pb.LogoutResponse{
		MessageLogout: fmt.Sprintf("User with ID %s successfully logged out", request.UserId),
	}, nil
}

func (repo *UserRepository) GetFollowersByUserID(ctx context.Context, request *pb.FollowersRequest) (*pb.FollowersResponse, error) {
	rows, err := repo.Db.QueryContext(ctx,
		`SELECT id, username, full_name FROM followers WHERE user_id = $1`,
		request.UserId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []*pb.Followers
	for rows.Next() {
		var id, username, fullName string
		if err := rows.Scan(&id, &username, &fullName); err != nil {
			return nil, err
		}
		follower := &pb.Followers{
			Id:       id,
			Username: username,
			FullName: fullName,
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.FollowersResponse{Followers: followers}, nil
}

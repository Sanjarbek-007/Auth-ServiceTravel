package postgres

import (
	pb "Auth-Service/genproto"
	storage "Auth-Service/help"
	"context"
	"database/sql"
	"fmt"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type UserRepository struct {
	Db *sql.DB
}


func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (repo *UserRepository) Register(request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    if repo.Db == nil {
        return nil, fmt.Errorf("database connection is not initialized")
    }

    var id, createdAt string
    err := repo.Db.QueryRow(
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

func (repo *UserRepository) Login(request *pb.LoginRequest) (*pb.LoginResponse, error) {
	var loginUser pb.LoginResponse
	err := repo.Db.QueryRow("select username, email, full_name from  users where email = $1 and password = $2", request.Email, request.Password).Scan(&loginUser.Username, &loginUser.Email, &loginUser.FullName)
	if err != nil {
		logge.Error("Error in Login")
		return nil, err
	}
	return &loginUser, nil
}

func (repo *UserRepository) Profile(request *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	var user pb.ProfileResponse
	err := repo.Db.QueryRow(
		"SELECT id, username, email, full_name, bio, countries_visited, created_at, updated_at from users WHERE id=$1 AND deleted_at IS NULL",
		request.UserId,
	).Scan(&user.Id, &user.Username, &user.Email, &user.FullName, &user.Bio, &user.CountriesVisited, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logge.Error("Error in Profile")
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateProfile(ctx context.Context, request *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	query := `UPDATE users 
			  SET full_name = $1, bio = $2, countries_visited = $3, updated_at = $4 
			  WHERE id = $5
			  RETURNING id, username, email, full_name, bio, countries_visited, updated_at`

	row := repo.Db.QueryRow(query,
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

func (repo *UserRepository) GetUsers(request *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	var (
		params = make(map[string]interface{})
		arr    []interface{}
		limit  string
		offset string
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
	query = query + filter + limit + offset
	query, arr = storage.ReplaceQueryParams(query, params)
	rows, err := repo.Db.Query(query, arr...)
	if err != nil {
		return nil, err
	}
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

func (repo *UserRepository) DeleteUser(request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, err := repo.Db.Exec(
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

func (repo *UserRepository) RefreshToken(ctx context.Context, request *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	refreshToken, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret_key), nil
	})

	if err != nil {
		logge.Error("error parsing refresh token", zap.Error(err))
		return nil, fmt.Errorf("error parsing refresh token: %v", err)
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {

		user := &pb.LoginResponse{
			Username: claims["username"].(string),
			FullName: claims["full_name"].(string),
			Email:    claims["email"].(string),
		}
		newTokens, err := repo.GENERATEJWTToken(user)
		if err != nil {
			return nil, fmt.Errorf("error generating new tokens: %v", err)
		}

		return &pb.RefreshResponse{
			AccessToken:  newTokens.AccessToken,
			RefreshToken: newTokens.RefreshToken,
		}, nil
	}

	return nil, fmt.Errorf("invalid refresh token")
}

func (repo *UserRepository) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	query := `UPDATE users
			SET token = NULL
			WHERE id = $1`

	result, err := repo.Db.Exec(query, request.UserId)
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
	rows, err := repo.Db.Query(
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


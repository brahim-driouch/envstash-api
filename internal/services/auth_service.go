package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
	"github.com/brahim-driouch/envstash.git/internal/utils"
	"github.com/brahim-driouch/envstash.git/internal/validators"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	authRepo interfaces.AuthRepository
}

func NewAuthService(r interfaces.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: r,
	}
}

func (s *AuthService) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	return s.authRepo.CreateRefreshToken(ctx, refreshToken)
}

func (s *AuthService) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return s.authRepo.FindRefreshToken(ctx, token)
}

func (s *AuthService) RevokeRefreshToken(ctx context.Context, token string) error {
	return s.authRepo.RevokeRefreshToken(ctx, token)
}

func (s *AuthService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.authRepo.RevokeAllUserTokens(ctx, userID)
}

func (s *AuthService) DeleteExpiredTokens(ctx context.Context) error {
	return s.authRepo.DeleteExpiredTokens(ctx)
}

func (s *AuthService) FindActiveUserTokens(ctx context.Context, userID string) (*[]models.RefreshToken, error) {
	return s.authRepo.FindActiveUserTokens(ctx, userID)
}

func (s *AuthService) DeleteUserToken(ctx context.Context, tokenID string, userID string) error {
	return s.authRepo.DeleteUserToken(ctx, tokenID, userID)
}
func (s *AuthService) RegisterUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	validationError := validators.ValidateNewUserFields(*input)
	if validationError != nil {
		return nil, validationError
	}

	userExists, err := s.authRepo.UserExists(ctx, input.Email)
	if err != nil {
		return nil, ErrUnexpected
	}
	if userExists {
		return nil, ErrUserExists
	}
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, ErrUnexpected
	}
	// Set the hashed password
	input.Password = string(hash)

	u, err := s.authRepo.CreateUser(ctx, input, input.Password)

	if err != nil {
		return nil, ErrUnexpected
	}
	return u, nil

}

func (s *AuthService) LoginUser(ctx context.Context, userLoginInput models.LoginInput) (*models.AuthToken, error) {
	// get the user from db
	user, err := s.authRepo.FindUserByEmail(ctx, userLoginInput.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	//if we have the user , compare passwords
	isValidPassword := utils.ComparePasswords(userLoginInput.Password, user.PasswordHash)
	if !isValidPassword {

		return nil, ErrInvalidCredentials
	}

	// genetate token
	var userSub = utils.TokenSub{
		Id:         user.ID,
		Fullname:   user.Fullname,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsAdmin:    user.IsAdmin,
	}
	// set access token err to 15 minutes
	accessToken, accessTokenErr := utils.GenerateAccessToken(userSub, 15)
	//set the refressh token for 30 dayas
	refreshToken, refreshTokenErr := utils.GenerateRefreshToken()

	if accessTokenErr != nil || refreshTokenErr != nil {
		return nil, ErrUnexpected
	}
	//store the refresh token in the database
	newRefreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(15 * 24 * time.Hour),
		CreatedAt: time.Now(),
		IPAddress: userLoginInput.IPAddress,
		UserAgent: userLoginInput.UserAgent,
	}
	err = s.authRepo.CreateRefreshToken(ctx, &newRefreshToken)
	if err != nil {
		log.Println("error creating refresh token", err)
		return nil, ErrUnexpected
	}
	return &models.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (s *AuthService) LogoutUser(ctx context.Context, token string) error {

	// get the refresh token from database
	refreshToken, err := s.authRepo.FindRefreshToken(ctx, token)
	if err != nil {
		return err
	}

	err = s.authRepo.RevokeRefreshToken(ctx, refreshToken.Token)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("no userID provided")
	}
	u, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

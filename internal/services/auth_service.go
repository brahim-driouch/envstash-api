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

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password too weak")
	ErrUserNotFound       = errors.New("no account found with this emails")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnexpected         = errors.New("could not proceed you request")
	ErrUserNotVerified    = errors.New("please verify your email before logging in. Check your inbox for the verification link.")
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

	go func(ctx context.Context) {

		verificationToken, err := utils.GenerateEmailVerifcationToken(u.ID)
		if err != nil {
			log.Printf("❌ Failed to generate verification token: %v", err)
			return
		}
		// Build email with template rendering
		emailParams, err := utils.BuildVerificationEmail(
			u.Email,
			u.Fullname,
			verificationToken,
		)

		if err != nil {
			log.Printf("❌ Failed to build verification email: %v", err)
			return
		}

		// Send email
		if err := utils.SendEmail(ctx, emailParams); err != nil {
			log.Printf("❌ Failed to send verification email to %s: %v", u.Email, err)
		} else {
			log.Printf("✓ Verification email sent to %s", u.Email)
		}
	}(ctx)
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
	//check if user is verified
	if !user.IsVerified {
		return nil, ErrUserNotVerified
	}

	// generate token
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

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("invalid verification token")
	}
	validToken, err := utils.VerifyAndDecodeEmailVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if validToken == nil {
		return errors.New("invalid verification token")
	}
	// check if user exists with the email from the token
	user, err := s.authRepo.FindUserByID(ctx, validToken.UserID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	if user.IsVerified {
		return errors.New("user already verified")
	}
	// update user as verified
	user.IsVerified = true
	err = s.authRepo.VerifyEmail(ctx, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if user.IsVerified {
		return errors.New("user already verified")
	}

	go func() {
		emailCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		verificationToken, err := utils.GenerateEmailVerifcationToken(user.ID)
		if err != nil {
			log.Printf("❌ Failed to generate verification token: %v", err)
			return
		}
		// Build email with template rendering
		emailParams, err := utils.BuildVerificationEmail(
			user.Email,
			user.Fullname,
			verificationToken,
		)

		if err != nil {
			log.Printf("❌ Failed to build verification email: %v", err)
			return
		}

		// Send email
		if err := utils.SendEmail(emailCtx, emailParams); err != nil {
			log.Printf("❌ Failed to send verification email to %s: %v", user.Email, err)
		} else {
			log.Printf("✓ Verification email sent to %s", user.Email)
		}

	}()
	return nil
}

package services

import (
	"context"
	"errors"
	"os"
	"time"
	"regexp"

	"FinalProject/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists        = errors.New("user with this email already exists")
	ErrWeakPassword      = errors.New("password does not meet complexity requirements")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

type AuthConfig struct {
	JWTSecret       []byte
	TokenExpiration time.Duration
	MaxLoginAttempts int
	RefreshTokenExpiration time.Duration
}

// AuthService manages authentication
type AuthService struct {
	UserRepo *repositories.UserRepository
	Config   AuthConfig
}

// NewAuthService initializes the service with configuration
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	config := AuthConfig{
		JWTSecret:       []byte(os.Getenv("JWT_SECRET")),
		TokenExpiration: 24 * time.Hour,
		MaxLoginAttempts: 5,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	}
	
	if len(config.JWTSecret) == 0 {
		panic("JWT_SECRET environment variable not set")
	}
	
	return &AuthService{
		UserRepo: userRepo,
		Config:   config,
	}
}

// validatePassword checks password complexity
func (s *AuthService) validatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	
	patterns := []string{
		`[A-Z]`,     // At least one uppercase letter
		`[a-z]`,     // At least one lowercase letter
		`[0-9]`,     // At least one number
		`[^A-Za-z0-9]`, // At least one special character
	}
	
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, password)
		if !matched {
			return ErrWeakPassword
		}
	}
	
	return nil
}

// HashPassword hashes the user password
func (s *AuthService) HashPassword(password string) (string, error) {
	if err := s.validatePassword(password); err != nil {
		return "", err
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// VerifyPassword checks if the password matches the stored hash
func (s *AuthService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// GenerateTokenPair generates both access and refresh tokens
func (s *AuthService) GenerateTokenPair(userID int, role string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateToken(userID, role, s.Config.TokenExpiration)
	if err != nil {
		return nil, err
	}
	
	// Generate refresh token
	refreshToken, err := s.generateToken(userID, role, s.Config.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(userID int, role string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(expiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.Config.JWTSecret)
}

// AuthenticateUser validates the email and password, then returns tokens
func (s *AuthService) AuthenticateUser(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.VerifyPassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	return s.GenerateTokenPair(user.ID, user.Role)
}

// ValidateToken validates and parses a JWT token
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.Config.JWTSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshTokens generates new token pair using a valid refresh token
func (s *AuthService) RefreshTokens(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID := int(claims["user_id"].(float64))
	role := claims["role"].(string)

	return s.GenerateTokenPair(userID, role)
}

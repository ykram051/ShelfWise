package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"FinalProject/models"
	"FinalProject/services"
)

type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController initializes the authentication controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

type RegisterInput struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required"`
	Role     string  `json:"role" validate:"required,oneof=admin customer"`
	Address  models.Address `json:"address"` 
}


type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Register handles user registration
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := c.AuthService.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Password does not meet complexity requirements", http.StatusBadRequest)
		return
	}

	// Assign default role if none is provided
	if input.Role == "" {
		input.Role = "customer"
	}

	// Create user
	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         input.Role,
		Address:      input.Address,  // ✅ Store Embedded Address
	}

	// Store user in database
	err = c.AuthService.UserRepo.CreateUser(ctx, &user)
	if err != nil {
		if err == services.ErrUserExists {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}


// Login handles user authentication
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	tokens, err := c.AuthService.AuthenticateUser(ctx, input.Email, input.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:   int64(c.AuthService.Config.TokenExpiration.Seconds()),
		TokenType:   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshToken handles token refresh requests
func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	tokens, err := c.AuthService.RefreshTokens(refreshToken)
	if err != nil {
		if err == services.ErrInvalidToken {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Token refresh failed", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:   int64(c.AuthService.Config.TokenExpiration.Seconds()),
		TokenType:   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

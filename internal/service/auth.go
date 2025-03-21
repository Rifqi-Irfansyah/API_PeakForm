package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/config"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type authService struct {
	cfg            *config.Config
	userRepository domain.UserRepository
}

func NewAuthService(cfg *config.Config, userRepository domain.UserRepository) domain.AuthService {
	return authService{cfg: cfg, userRepository: userRepository}
}

// Login authenticates a user by verifying the provided credentials and generates a JWT token upon success.
func (a authService) Login(ctx context.Context, data dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.userRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		log.Println("Error while searching for user:", err)
		return dto.AuthResponse{}, err
	}

	if user.ID == "" {
		log.Println("User not found or invalid email:", data.Email)
		return dto.AuthResponse{}, errors.New("username or password is wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		log.Println("Password does not match for user:", data.Email)
		return dto.AuthResponse{}, errors.New("username or password is wrong")
	}

	tokenStr, err := generateJWT(user.ID, a.cfg.Jwt.Key, a.cfg.Jwt.Exp)
	if err != nil {
		log.Println("Failed to generate JWT for user:", data.Email, "error:", err)
		return dto.AuthResponse{}, errors.New("failed to generate token")
	}

	log.Println("Login successful for user:", data.Email)
	return dto.AuthResponse{Token: tokenStr}, nil
}

// Register creates a new user account, hashes the password, saves the user in the repository, and returns a JWT token.
func (a authService) Register(ctx context.Context, data dto.RegisterRequest) error {
	if data.Name == "" || data.Email == "" || data.Password == "" {
		log.Println("Register error: missing required fields")
		return errors.New("all fields are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Register error: failed to hash password:", err)
		return errors.New("failed to hash password")
	}

	user := domain.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(hashedPassword),
	}

	err = a.userRepository.Save(ctx, user)
	if err != nil {
		log.Println("Register error: failed to save user:", err)
		return errors.New("failed to save user")
	}

	return nil
}

// generateJWT generates a JWT token string for a given user ID, secret key, and expiration time in minutes.
// It returns the signed token string or an error if signing fails.
func generateJWT(userID string, secretKey string, expMinutes int) (string, error) {
	if secretKey == "" {
		log.Println("Secret key is empty!")
		return "", errors.New("secret key is empty")
	}

	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

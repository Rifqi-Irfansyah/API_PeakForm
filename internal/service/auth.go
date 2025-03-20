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
		log.Println("Error saat mencari user:", err)
		return dto.AuthResponse{}, err
	}

	if user.ID == "" {
		log.Println("User tidak ditemukan atau email salah:", data.Email)
		return dto.AuthResponse{}, errors.New("username or password is wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		log.Println("Password tidak cocok untuk user:", data.Email)
		return dto.AuthResponse{}, errors.New("username or password is wrong")
	}

	tokenStr, err := generateJWT(user.ID, a.cfg.Jwt.Key, a.cfg.Jwt.Exp)
	if err != nil {
		log.Println("Gagal generate JWT untuk user:", data.Email, "error:", err)
		return dto.AuthResponse{}, errors.New("failed to generate token")
	}

	log.Println("Login berhasil untuk user:", data.Email)
	return dto.AuthResponse{Token: tokenStr}, nil
}

// generateJWT generates a JWT token string for a given user ID, secret key, and expiration time in minutes.
// It returns the signed token string or an error if signing fails.
func generateJWT(userID string, secretKey string, expMinutes int) (string, error) {
	if secretKey == "" {
		log.Println("Secret key kosong!")
		return "", errors.New("secret key is empty")
	}

	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"time"
)

type authService struct {
	cfg            *config.Config
	userRepository domain.UserRepository
	otpRepository  domain.OtpRepository
}

func NewAuthService(cfg *config.Config, userRepository domain.UserRepository, otpRepository domain.OtpRepository) domain.AuthService {
	return authService{cfg: cfg, userRepository: userRepository, otpRepository: otpRepository}
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

func (a authService) ForgotPassword(ctx context.Context, email string) error {
	// Cek apakah email user ada
	user, err := a.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("failed to check user")
	}
	if user.ID == "" {
		return errors.New("email not found")
	}

	otp := generateOTP()

	err = a.otpRepository.SaveOTP(ctx, email, otp, time.Now().Add(10*time.Minute))
	if err != nil {
		return errors.New("failed to save OTP")
	}

	err = sendOTPByEmail(email, otp)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (a authService) ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error {
	valid, err := a.otpRepository.VerifyOTP(ctx, data.Email, data.OTP)
	if err != nil || !valid {
		return errors.New("invalid or expired OTP")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = a.userRepository.UpdatePassword(ctx, data.Email, string(hashedPassword))
	if err != nil {
		return errors.New("failed to update password")
	}

	err = a.otpRepository.DeleteOTP(ctx, data.Email)
	if err != nil {
		return errors.New("failed to delete OTP")
	}

	return nil
}

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func sendOTPByEmail(email, otp string) error {
	from := "fznrzkxxviii@gmail.com"
	password := "jxkg vpen cizu xsap"
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	message := fmt.Sprintf(`
	Subject: 🔐 Reset Password - OTP Verification
	
	MIME-version: 1.0;
	Content-Type: text/html; charset="UTF-8";
	
	<html>
	<head>
		<style>
			.container {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				color: #333;
				text-align: center;
			}
			.otp {
				font-size: 24px;
				font-weight: bold;
				color: #007bff;
			}
			.footer {
				margin-top: 20px;
				font-size: 12px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>🔒 OTP Verification for Password Reset</h2>
			<p>Hello,</p>
			<p>We received a request to reset your account password. Use the following OTP code to proceed with the password reset process:</p>
			<p class="otp">%s</p>
			<p>Do not share this code with anyone for the security of your account.</p>
			<p>If you did not request a password reset, please ignore this email.</p>
			<p class="footer">© 2025 PeakForm. All Rights Reserved.</p>
		</div>
	</body>
	</html>`, otp)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "🔐 Your OTP Code for Password Reset")
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)

	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("OTP sent successfully to", email)
	return nil
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

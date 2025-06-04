package service

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"api-peak-form/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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

func (a authService) Login(ctx context.Context, data dto.AuthRequest) (dto.AuthResponse, error) {
	logrus.Infof("Attempting login for user: %s", data.Email)
	user, err := a.userRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		logrus.Errorf("Error while searching for user: %v", err)
		return dto.AuthResponse{}, err
	}

	if user.ID == "" {
		logrus.Warnf("User not found or invalid email: %s", data.Email)
		return dto.AuthResponse{}, errors.New("username or password is wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		logrus.Warnf("Password does not match for user: %s", data.Email)
		return dto.AuthResponse{}, errors.New("username or password is wrong")
	}

	tokenStr, err := generateJWT(user.ID, a.cfg.Jwt.Key, a.cfg.Jwt.Exp)
	if err != nil {
		logrus.Errorf("Failed to generate JWT for user: %s, error: %v", data.Email, err)
		return dto.AuthResponse{}, errors.New("failed to generate token")
	}

	logrus.Infof("Login successful for user: %s", data.Email)
	return dto.AuthResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Point:  user.Point,
		Streak: user.Streak,
		Token:  tokenStr,
	}, nil
}

func (a authService) Register(ctx context.Context, data dto.RegisterRequest) error {
	if data.Email == "" {
		logrus.Warn("Register error: missing required fields")
		return errors.New("all fields are required")
	}

	existingUser, err := a.userRepository.FindByEmail(ctx, data.Email)
	if err == nil && existingUser.ID != "" {
		logrus.Warnf("Email already registered: %s", data.Email)
		return errors.New("email already registered")
	}

	otp := generateOTP()
	logrus.Infof("Generated OTP for email %s: %s", data.Email, otp)

	err = a.otpRepository.SaveOTP(ctx, data.Email, otp, time.Now().Add(10*time.Minute))
	if err != nil {
		logrus.Errorf("Failed to save OTP for email %s: %v", data.Email, err)
		return errors.New("failed to save OTP")
	}

	err = sendOTPByEmail(data.Email, otp)
	if err != nil {
		logrus.Errorf("Failed to send OTP to email %s: %v", data.Email, err)
		return errors.New("failed to send OTP")
	}

	logrus.Infof("OTP sent successfully to email: %s", data.Email)
	return nil
}

func (a authService) VerifyRegisterOTP(ctx context.Context, data dto.VerifyOTPRequest) error {
	logrus.Infof("Verifying OTP for email: %s", data.Email)

	valid, err := a.otpRepository.VerifyOTP(ctx, data.Email, data.OTP)
	if err != nil || !valid {
		logrus.Warnf("Invalid or expired OTP for email: %s", data.Email)
		return errors.New("invalid or expired OTP")
	}

	err = a.otpRepository.DeleteOTP(ctx, data.Email)
	if err != nil {
		logrus.Errorf("Failed to delete OTP for email: %s", data.Email)
		return errors.New("failed to delete OTP")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("Failed to hash password:", err)
		return errors.New("failed to hash password")
	}

	user := domain.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(hashedPassword),
	}

	err = a.userRepository.Save(ctx, user)
	if err != nil {
		logrus.Errorf("Failed to save user: %v", err)
		return errors.New("failed to save user")
	}

	logrus.Infof("User registered successfully: %s", data.Email)
	return nil
}

func (a authService) ForgotPassword(ctx context.Context, email string) error {
	logrus.Infof("Attempting to process forgot password for email: %s", email)
	user, err := a.userRepository.FindByEmail(ctx, email)
	if err != nil {
		logrus.Errorf("Failed to check user for email %s: %v", email, err)
		return errors.New("failed to check user")
	}
	if user.ID == "" {
		logrus.Warnf("Email not found: %s", email)
		return errors.New("email not found")
	}

	otp := generateOTP()
	logrus.Infof("Generated OTP for email %s: %s", email, otp)

	err = a.otpRepository.SaveOTP(ctx, email, otp, time.Now().Add(10*time.Minute))
	if err != nil {
		logrus.Errorf("Failed to save OTP for email %s: %v", email, err)
		return errors.New("failed to save OTP")
	}

	err = sendOTPByEmail(email, otp)
	if err != nil {
		logrus.Errorf("Failed to send OTP email to %s: %v", email, err)
		return errors.New(err.Error())
	}

	logrus.Infof("OTP sent successfully to email: %s", email)
	return nil
}

func (a authService) ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error {
	logrus.Infof("Attempting to reset password for email: %s", data.Email)
	valid, err := a.otpRepository.VerifyOTP(ctx, data.Email, data.OTP)
	if err != nil || !valid {
		logrus.Warnf("Invalid or expired OTP for email: %s", data.Email)
		return errors.New("invalid or expired OTP")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed to hash new password for email: %s, error: %v", data.Email, err)
		return errors.New("failed to hash password")
	}

	err = a.userRepository.UpdatePassword(ctx, data.Email, string(hashedPassword))
	if err != nil {
		logrus.Errorf("Failed to update password for email: %s, error: %v", data.Email, err)
		return errors.New("failed to update password")
	}

	err = a.otpRepository.DeleteOTP(ctx, data.Email)
	if err != nil {
		logrus.Errorf("Failed to delete OTP for email: %s, error: %v", data.Email, err)
		return errors.New("failed to delete OTP")
	}

	logrus.Infof("Password reset successfully for email: %s", data.Email)
	return nil
}

func generateOTP() string {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	logrus.Infof("Generated OTP: %s", otp)
	return otp
}

func sendOTPByEmail(email, otp string) error {
	from := "fznrzkxxviii@gmail.com"
	password := "jxkg vpen cizu xsap"
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	message := fmt.Sprintf(`
	
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
			<h2>🔒 OTP Verification</h2>
			<p>Hello,</p>
			<p>Use the following OTP code to proceed with the verification process:</p>
			<p class="otp">%s</p>
			<p>Do not share this code with anyone for the security of your account.</p>
			<p>If you did not request this, please ignore this email.</p>
			<p class="footer">© 2025 PeakForm. All Rights Reserved.</p>
		</div>
	</body>
	</html>`, otp)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "🔐 Your OTP Code for Verification")
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)

	if err := dialer.DialAndSend(mailer); err != nil {
		logrus.Errorf("Error sending email to %s: %v", email, err)
		return err
	}

	logrus.Infof("OTP sent successfully to %s", email)
	return nil
}

func generateJWT(userID string, secretKey string, expMinutes int) (string, error) {
	if secretKey == "" {
		logrus.Error("Secret key is empty")
		return "", errors.New("secret key is empty")
	}

	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logrus.Errorf("Failed to sign token: %v", err)
		return "", err
	}

	logrus.Infof("JWT generated successfully for user ID: %s", userID)
	return tokenStr, nil
}

func (a authService) ChangePassword(ctx context.Context, data dto.ChangePasswordRequest) error {
	logrus.Infof("Attempting to change password for user ID: %s", data.ID)

	user, err := a.userRepository.FindByID(ctx, data.ID)
	if err != nil {
		logrus.Errorf("Failed to find user: %v", err)
		return errors.New("failed to find user")
	}

	if user.ID == "" {
		logrus.Warnf("User not found for ID: %s", data.ID)
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword)); err != nil {
		logrus.Warnf("Old password does not match for user ID: %s", data.ID)
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed to hash new password for user ID: %s, error: %v", data.ID, err)
		return errors.New("failed to hash new password")
	}

	err = a.userRepository.UpdatePassword(ctx, user.Email, string(hashedPassword))
	if err != nil {
		logrus.Errorf("Failed to update password for user ID: %s, error: %v", data.ID, err)
		return errors.New("failed to update password")
	}

	logrus.Infof("Password changed successfully for user ID: %s", data.ID)
	return nil
}

func (a authService) CheckToken(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.cfg.Jwt.Key), nil
	})

	if err != nil || !parsedToken.Valid {
		logrus.Warnf("Invalid token: %v", err)
		return errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.Error("Failed to parse token claims")
		return errors.New("failed to parse token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		logrus.Warn("Token has expired")
		return errors.New("token has expired")
	}

	logrus.Infof("Token is valid and not expired")
	return nil
}

func (a authService) GetUserByToken(ctx context.Context, token string) (dto.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.cfg.Jwt.Key), nil
	})

	if err != nil || !parsedToken.Valid {
		logrus.Warnf("Invalid token: %v", err)
		return dto.AuthResponse{}, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.Error("Failed to parse token claims")
		return dto.AuthResponse{}, errors.New("failed to parse token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		logrus.Warn("Token has expired")
		return dto.AuthResponse{}, errors.New("token has expired")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		logrus.Error("User ID not found in token claims")
		return dto.AuthResponse{}, errors.New("user ID not found in token claims")
	}

	user, err := a.userRepository.FindByID(ctx, userID)
	if err != nil {
		logrus.Errorf("Failed to find user by ID: %v", err)
		return dto.AuthResponse{}, errors.New("failed to find user")
	}

	logrus.Infof("User retrieved successfully from token: %s", userID)
	return dto.AuthResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Point:  user.Point,
		Streak: user.Streak,
	}, nil
}

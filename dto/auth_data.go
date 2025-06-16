package dto

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Point    int    `json:"point"`
	Streak   int    `json:"streak"`
	PhotoURL string `json:"photo"`
	Rank     int    `json:"rank"`
	Token    string `json:"token"`
}

type RegisterRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type VerifyOTPRequest struct {
	Email    string `json:"email" validate:"required,email"`
	OTP      string `json:"otp" validate:"required,len=6"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OTP         string `json:"otp" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ChangePasswordRequest struct {
	ID          string `json:"id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

package dto

type UserStreakRequest struct {
	UserID string `json:"id" validate:"required"`
	Streak bool   `json:"streak" validate:"required"`
}

type UserLeaderboardResponse struct {
	Name  string `json:"name"`
	Point int    `json:"point"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Point    int    `json:"point"`
	Streak   int    `json:"streak"`
	PhotoURL string `json:"photo"`
}

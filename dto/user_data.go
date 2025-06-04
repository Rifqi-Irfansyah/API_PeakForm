package dto

type UserStreakRequest struct {
	UserID string `json:"id" validate:"required"`
	Streak bool   `json:"streak" validate:"required"`
}

type UserLeaderboardResponse struct {
	Name  string `json:"name"`
	Point int    `json:"point"`
}

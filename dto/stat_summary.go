package dto

type StatsSummary struct {
	UserID             string            `json:"user_id"`
	TotalSets          int               `json:"total_sets"`
	TotalReps          int               `json:"total_repetitions"`
	ExerciseCounter    map[string]int    `json:"exercise_counter"`
	TotalExerciseCount int               `json:"total_exercise_count"`
}
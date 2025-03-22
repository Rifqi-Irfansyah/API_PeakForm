package dto

type ScheduleRequest struct {
	UID string `json:"id"`
}

type ScheduleListResponse struct {
	Schedules []ScheduleResponse `json:"schedules"`
}

type ExerciseResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Muscle       string `json:"muscle"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	Instructions string `json:"instructions"`
	Gif          string `json:"gif"`
}

type ScheduleResponse struct {
	ID        uint      `json:"id"`
	Day       int         `json:"day"`
	Exercises []ExerciseResponse 	  `json:"exercise"`
}
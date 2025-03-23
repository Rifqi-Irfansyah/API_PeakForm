package dto

type CreateScheduleRequest struct{
	UID 		string 	`json:"id_user" form:"id_user"`
	ExerciseID 	uint	`json:"id_exercise" form:"id_exercise"`
	Type		string	`json:"type"  form:"type"`
	Day			int 	`json:"day"`
	Set			int		`json:"set"`
	Repetition 	int 	`json:"repetition"` 	
}

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
	ID        	uint      	`json:"id"`
	Day       	int         `json:"day"`
	Type		string		`json:"type"`
	Exercises []ExerciseResponse 	  `json:"exercise"`
}
package dto

type CreateScheduleRequest struct{
	UID 		string 	`json:"id_user" form:"id_user"`
	ExerciseID 	uint	`json:"id_exercise" form:"id_exercise"`
	Type		string	`json:"type"  form:"type"`
	Day			int 	`json:"day"`
	Set			int		`json:"set"`
	Repetition 	int 	`json:"repetition"` 	
}

type UpdateScheduleRequest struct{
	ID			string	`json:"id" form:"id"`
	Day			int 	`json:"day" form:"day"`
}

type UpdateExerciseScheduleRequest struct{
	ID			string	`json:"id" form:"id"`
	ExerciseID 	uint	`json:"id_exercise" form:"id_exercise"`
	NExerciseID uint	`json:"new_id_exercise" form:"new_id_exercise"`
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
	Set 		 int	`json:"set"`
	Repetition 	 int	`json:"repetition"`
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
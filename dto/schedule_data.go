package dto

type CreateScheduleRequest struct {
	UID        string `json:"id_user" form:"id_user" validate:"required"`
	ExerciseID uint   `json:"id_exercise" form:"id_exercise" validate:"required"`
	Type       string `json:"type"  form:"type" validate:"required"`
	Day        int    `json:"day" validate:"required,gt=0,lte=7"`
	Set        int    `json:"set" validate:"required,gt=0"`
	Repetition int    `json:"repetition" validate:"required,gt=0"`
}

type UpdateScheduleRequest struct {
	ID  string `json:"id" form:"id" validate:"required"`
	Day int    `json:"day" form:"day" validate:"gt=0,lte=7"`
}

type UpdateExerciseScheduleRequest struct {
	ID          string `json:"id" form:"id" validate:"required"`
	ExerciseID  uint   `json:"id_exercise" form:"id_exercise"`
	NExerciseID uint   `json:"new_id_exercise" form:"new_id_exercise"`
	Set         int    `json:"set"`
	Repetition  int    `json:"repetition"`
}

type ScheduleRequest struct {
	UID string `json:"id" validate:"required"`
}

type ScheduleListResponse struct {
	Schedules []ScheduleResponse `json:"schedules"`
}

type ExerciseResponse struct {
	ID           uint   `json:"id"`
	Set          int    `json:"set"`
	Repetition   int    `json:"repetition"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Muscle       string `json:"muscle"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	Instructions string `json:"instructions"`
	Image        string `json:"image"`
}

type ScheduleResponse struct {
	ID        uint               `json:"id"`
	Day       int                `json:"day"`
	Type      string             `json:"type"`
	Exercises []ExerciseResponse `json:"exercise"`
}

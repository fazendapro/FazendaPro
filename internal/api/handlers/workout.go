package handlers

type WorkoutHandler struct {
}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
}

func (h *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
}
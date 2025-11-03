package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Naveenravi07/go-api/internal/store"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.Error(w, "Id must be an int", http.StatusBadRequest)
		return
	}

	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	fmt.Fprintf(w, "%+v", workout)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Printf("failed to decode workout request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Printf("\n\nfailed to create new workout : %+v \n\n", err)
		http.Error(w, "failed to create a workout", http.StatusBadGateway)
		return
	}

	fmt.Fprintf(w, "New workout created with id : %d \n", createdWorkout.Id)
}

package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Naveenravi07/go-api/internal/store"
	"github.com/Naveenravi07/go-api/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v ", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id "})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutById: %v ", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout did not exist "})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"data": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: invalid req.body %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: failed to create workout %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"data": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: invalid req.body %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}

	err = wh.workoutStore.UpdateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: failed to update workout %v", err)
	}
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"data": "workout updated successfully"})
}

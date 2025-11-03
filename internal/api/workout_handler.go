package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {

}

func  NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
}

func (wh *WorkoutHandler) HandleWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId,err := strconv.ParseInt(paramsWorkoutId,10,64)
	if err != nil{
		http.Error(w,"Id must be an int",http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w,"This is the workout id %d",workoutId)
}


func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"Created workout \n");
}

package routes

import (
	"github.com/Naveenravi07/go-api/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleWorkoutById)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Patch("/workouts", app.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workouts/{id}", app.WorkoutHandler.DeleteWorkoutHandler)

	r.Post("/user", app.UserHandler.CreateUserHandler)
	r.Patch("/user", app.UserHandler.UpdateUserHandler)
	r.Get("/user/{username}", app.UserHandler.GetUserByUsernameHandler)

	return r
}

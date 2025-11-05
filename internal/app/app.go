package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Naveenravi07/go-api/internal/api"
	"github.com/Naveenravi07/go-api/internal/store"
	"github.com/Naveenravi07/go-api/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	userStore := store.NewPostgresUserStore(pgDB)
	userHander := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:         logger,
		DB:             pgDB,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHander,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running \n")
}

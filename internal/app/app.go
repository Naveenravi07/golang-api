package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Naveenravi07/go-api/internal/api"
	"github.com/Naveenravi07/go-api/internal/store"
)

type Application struct {
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pgDB,err := store.Open()
	if err != nil{
		return  nil,err
	}

	app := &Application{
		Logger: logger,
		DB: pgDB,
		WorkoutHandler: api.NewWorkoutHandler(),
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running \n")
}

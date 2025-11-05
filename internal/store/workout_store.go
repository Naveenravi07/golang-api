package store

import (
	"database/sql"
	"errors"
	"fmt"
)

type Workout struct {
	Id              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	Id              int      `json:"id"`
	WorkoutId       int      `json:"workout_id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float32 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutById(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query :=
		`INSERT INTO workouts (title,description,duration_minutes,calories_burned)
	VALUES ($1,$2,$3,$4)
	RETURNING id;`

	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.Id)
	if err != nil {
		return nil, err
	}

	for i := range workout.Entries {
		entry := &workout.Entries[i]
		query := `
		INSERT INTO workout_entries (workout_id,exercise_name,sets,reps,duration_seconds,weight,notes,order_index)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id;
		`
		err = tx.QueryRow(query, workout.Id, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.Id)
		if err != nil {
			return nil, err
		}
		entry.WorkoutId = workout.Id
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `SELECT id,title,duration_minutes,calories_burned from workouts where id=$1`

	err := pg.db.QueryRow(query, id).Scan(&workout.Id, &workout.Title, &workout.DurationMinutes, &workout.CaloriesBurned)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	query = `SELECT id,workout_id,exercise_name,sets,reps,duration_seconds,weight,notes,order_index from workout_entries where workout_id=$1 ORDER BY order_index`
	results, err := pg.db.Query(query, workout.Id)
	if err != nil {
		return nil, err
	}

	defer results.Close()
	for results.Next() {
		var entry WorkoutEntry
		err = results.Scan(&entry.Id, &entry.WorkoutId, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	query := `UPDATE workouts set title=$1,description=$2,duration_minutes=$3,calories_burned=$4 where id=$5`
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	rows, err := tx.Query(`SELECT id,workout_id,exercise_name,sets,reps,duration_seconds,weight,notes,order_index FROM workout_entries where workout_id=$1`, workout.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	currentIds := map[int]bool{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		currentIds[id] = true
	}

	newIds := map[int]bool{}
	for _, entry := range workout.Entries {
		if entry.Id == 0 {
			insertQ := `
				INSERT INTO workout_entries 
				(workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			`
			_, err := tx.Exec(insertQ, workout.Id, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex)
			if err != nil {
				return err
			}
		} else {
			updateQ := `
				UPDATE workout_entries
				SET exercise_name=$1, sets=$2, reps=$3, duration_seconds=$4, weight=$5, notes=$6, order_index=$7
				WHERE id=$8 AND workout_id=$9
			`
			_, err := tx.Exec(updateQ, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex, entry.Id, workout.Id)
			if err != nil {
				return err
			}
			newIds[entry.Id] = true
		}
	}

	for id := range newIds {
		if !currentIds[id] {
			_, err := tx.Exec(`DELETE FROM workout_entries WHERE id=$1`, id)
			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `DELETE FROM workouts where id=$1`
	result, err := pg.db.Exec(query, id)
	fmt.Printf("\n\n %+v \n\n", result)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No workout found for id ")
	}
	return nil
}

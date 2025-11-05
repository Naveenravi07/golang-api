package store

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db : %v", err)
	}

	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Migrating test db error : %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts,workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("Error on cleaning up db before test : %v", err)
	}

	return db
}

func TestWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	lastInsertedId := 0
	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Push day",
				Description:     "Push harder",
				DurationMinutes: 12,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						Id:              1,
						WorkoutId:       1,
						ExerciseName:    "Push up",
						Sets:            2,
						Reps:            IntPtr(3),
						DurationSeconds: nil,
						Weight:          FloatPtr(10),
						Notes:           "Good had some more energy lefttt",
						OrderIndex:      1,
					},
					{
						Id:              2,
						WorkoutId:       1,
						ExerciseName:    "Pull up",
						Sets:            2,
						Reps:            IntPtr(3),
						DurationSeconds: nil,
						Weight:          FloatPtr(10),
						Notes:           "Not perfect but okay",
						OrderIndex:      2,
					},
					{
						Id:              3,
						WorkoutId:       1,
						ExerciseName:    "Thread mill run",
						Sets:            3,
						Reps:            nil,
						DurationSeconds: IntPtr(4000),
						Weight:          nil,
						Notes:           "ok",
						OrderIndex:      3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid workout",
			workout: &Workout{
				Title:           "Full body ",
				Description:     "Complete workout",
				DurationMinutes: 12,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						Id:              1,
						WorkoutId:       1,
						ExerciseName:    "Bench press",
						Sets:            2,
						Reps:            IntPtr(3),
						DurationSeconds: IntPtr(300),
						Weight:          FloatPtr(10),
						Notes:           "Good had some more energy lefttt",
						OrderIndex:      1,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrieved, err := store.GetWorkoutById(int64(tt.workout.Id))
			fmt.Printf("\n\n %+v \n\n", retrieved)
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.Id, retrieved.Id)
			assert.Equal(t, createdWorkout.Title, retrieved.Title)
			assert.Equal(t, len(createdWorkout.Entries), len(tt.workout.Entries))

			for i, _ := range retrieved.Entries {
				assert.Equal(t, retrieved.Entries[i].ExerciseName, tt.workout.Entries[i].ExerciseName)
				assert.Equal(t, retrieved.Entries[i].Id, createdWorkout.Entries[i].Id)
				assert.Equal(t, retrieved.Entries[i].Sets, tt.workout.Entries[i].Sets)
			}
			lastInsertedId = createdWorkout.Id
		})
	}

	deleteTests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Delete workout valid",
			id:      lastInsertedId,
			wantErr: false,
		},
		{
			name:    "Delete workout invalid id ",
			id:      2,
			wantErr: true,
		},
	}

	for _, tt := range deleteTests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.DeleteWorkout(int64(tt.id))
			if tt.wantErr {
				fmt.Printf("\n\n Deletion error : %+v \n\n", err)
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			workout, err := store.GetWorkoutById(int64(tt.id))
			if workout != nil {
				t.Errorf("Expected workout to be deleted, but found %+v", workout)
			}
			if err == nil {
				t.Errorf("Expected error when retrieving deleted workout, but got none")
			}
		})
	}
}

func FloatPtr(v float32) *float32 { return &v }
func IntPtr(v int) *int           { return &v }

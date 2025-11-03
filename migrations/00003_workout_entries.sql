-- +goose Up
-- +goose statementBegin
CREATE TABLE IF NOT EXISTS workout_entries(
    id BIGSERIAL PRIMARY KEY,
    workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_name varchar(255) NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER,
    duration_seconds INTEGER,
    weight DECIMAL(5,2),
    notes TEXT,
    order_index INTEGER NOT NULL,
    createdAT TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_workout_entry CHECK(
        (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND 
        (reps IS NULL OR duration_seconds IS NULL)
    )
);
-- +goose statementEnd
-- +goose Down

-- +goose statementBegin
DROP TABLE workout_entries;
-- +goose statementEnd



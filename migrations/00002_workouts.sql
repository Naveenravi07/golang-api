-- +goose Up
-- +goose statementBegin
CREATE TABLE IF NOT EXISTS workouts(
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL,
    calories_burned INTEGER,
    createdAT TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose statementEnd
-- +goose Down

-- +goose statementBegin
DROP TABLE workouts;
-- +goose statementEnd


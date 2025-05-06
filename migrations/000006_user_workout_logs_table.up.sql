CREATE TABLE user_workout_logs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    workout_name TEXT NOT NULL,
    sets INT,
    reps INT,
    duration_minutes INT,
    weight_kg INT,
    log_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ DEFAULT now()
);

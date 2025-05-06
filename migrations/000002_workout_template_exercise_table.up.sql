CREATE TABLE IF NOT EXISTS workout_exercises (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    template_id INT REFERENCES workout_templates(id),
    name TEXT NOT NULL,         
    sets INT,
    reps INT
);

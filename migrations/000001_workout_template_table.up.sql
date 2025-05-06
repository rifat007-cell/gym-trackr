CREATE TABLE IF NOT EXISTS workout_templates (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    goal TEXT NOT NULL,         
    level TEXT NOT NULL,        
    name TEXT NOT NULL,
    description TEXT NOT NULL
);



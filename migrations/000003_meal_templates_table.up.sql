CREATE TABLE IF NOT EXISTS meal_templates (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    goal TEXT NOT NULL,         
    dietary_preference TEXT,    
    name TEXT NOT NULL,         
    description TEXT,
    calories INT
);

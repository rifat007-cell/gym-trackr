CREATE TABLE user_daily_plan (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    date TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    workout_template_id BIGINT REFERENCES workout_templates(id),
    meal_template_id BIGINT REFERENCES meal_templates(id) 
);

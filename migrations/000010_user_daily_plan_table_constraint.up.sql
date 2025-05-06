ALTER TABLE user_daily_plan
ADD CONSTRAINT unique_user_date UNIQUE (user_id, date);
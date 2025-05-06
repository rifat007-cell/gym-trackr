package data

import (
	"context"
	"database/sql"
	"time"
)

type UserDailyPlan struct{
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Date        time.Time `json:"date"`
	WorkoutTemplateID int  `json:workout_templateid`

	MealTemplateID int  `json:meal_templateid`

}

type UserDailyPlanModel struct{
	DB *sql.DB
}

func (u *UserDailyPlanModel) InsertDailyPlan(userID int, workoutTemplateID int) (int, error) {
	stmt:= `INSERT INTO user_daily_plan (user_id, date, workout_template_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, date) DO UPDATE 
		SET workout_template_id = EXCLUDED.workout_template_id
		RETURNING id

	`

	var id int

	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)

	defer cancel()

	// Use only date part
	today := time.Now().Truncate(24 * time.Hour)

 	
	err:=u.DB.QueryRowContext(ctx,stmt,userID,today,workoutTemplateID).Scan(&id)
	if err!=nil{
		return 0,err
	}

	return id,nil





}


func (u *UserDailyPlanModel) UpdateDailyPlan(mealTemplateID,userID int) (int, error) {
	stmt:= `INSERT INTO user_daily_plan (user_id, date, meal_template_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, date) DO UPDATE 
		SET meal_template_id = EXCLUDED.meal_template_id
		RETURNING id

	`


	var id int

	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)

	defer cancel()

	// Use only date part
	today := time.Now().Truncate(24 * time.Hour)

 	
	err:=u.DB.QueryRowContext(ctx,stmt,userID,today,mealTemplateID).Scan(&id)
	if err!=nil{
		return 0,err
	}

	return id,nil

}


func (u *UserDailyPlanModel) GetDailyPlanByUserID(userID int) (UserDailyPlan, error) {
	stmt := `SELECT id, user_id, date, workout_template_id, meal_template_id
	FROM user_daily_plan
	WHERE user_id = $1`

	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)

	defer cancel()

	var dailyPlan UserDailyPlan

	err := u.DB.QueryRowContext(ctx, stmt, userID).Scan(&dailyPlan.ID, &dailyPlan.UserID, &dailyPlan.Date, &dailyPlan.WorkoutTemplateID, &dailyPlan.MealTemplateID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return UserDailyPlan{} , nil
		}
		return UserDailyPlan{}, err 
	}
	return dailyPlan, nil
	
}


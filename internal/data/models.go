package data

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
)
var (
    ErrRecordNotFound = errors.New("record not found")
		ErrEditConflict   = errors.New("edit conflict")
		ErrInvalidCredentials = errors.New("invalid credentials")

)

type Models struct{
  Workouts WorkoutModel
	Exercises ExerciseModel
	Meals     MealsModel
	Users     UserModel
	Tokens    TokenModel
	WorkoutLogs WorkoutLogModel
	Passkey     PasskeyModel
	UserDailyPlan UserDailyPlanModel
	ChatMessages ChatMessageModel
	
}

func NewModels(db *sql.DB) Models {
	return Models{
		Workouts: WorkoutModel{DB: db},
		Exercises: ExerciseModel{DB: db},
		Meals:     MealsModel{DB: db},
		Users:     UserModel{DB: db},
		Tokens: 	TokenModel{DB: db},
		WorkoutLogs: WorkoutLogModel{DB: db},
		Passkey:     *NewPasskeyModel(db,slog.New(slog.NewTextHandler(os.Stdout, nil))),
		UserDailyPlan: UserDailyPlanModel{DB: db},
		ChatMessages: ChatMessageModel{DB: db},
	}
}
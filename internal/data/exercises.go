package data

import "database/sql"


type Exercise struct{
	ID  int `json:"id"`
	Name string `json:"name"`
	Sets int `json:"sets"`
	Reps int `json:"reps"`
}

type ExerciseModel struct{
	DB *sql.DB
}
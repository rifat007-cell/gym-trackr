package data

import (
	"context"
	"database/sql"
	"time"
)


type Workout struct {
	ID          int    `json:"id"`
	Goal        string `json:"goal"`
	Level       string  `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Exercises   []Exercise `json:"exercises"`
}


type WorkoutModel struct{
	DB *sql.DB
}

func (m WorkoutModel) GetAllExerciseBasedWorkoutName(goal,level string)([]Workout,error){
	// this query is called the "subquery join"
	stmt:= `SELECT w.id, w.name,w.goal,w.level,w.description,we.id,we.name,we.sets,we.reps
	        FROM (
					   SELECT w.id, w.name,w.goal,w.level,w.description
						 FROM workout_templates AS w
						 WHERE w.goal = $1 AND w.level = $2
						 ORDER BY RANDOM() LIMIT 1
					)
					 w 
					JOIN workout_exercises AS we ON w.id = we.template_id
	`


	ctx,cancel:= context.WithTimeout(context.Background(),3*time.Second)

	defer cancel()

	rows,err:= m.DB.QueryContext(ctx,stmt,goal,level)

	if err!=nil{
		return nil,err
	}
	defer rows.Close()
  
	var workouts []Workout

	for rows.Next(){
		var w Workout
		var e Exercise

		err= rows.Scan(&w.ID,&w.Name,&w.Goal,&w.Level,&w.Description,&e.ID,&e.Name,&e.Sets,&e.Reps)
		if err!=nil{
			return nil,err
		}

		w.Exercises= append(w.Exercises,e)

		workouts= append(workouts,w)
	}
	if err= rows.Err(); err!=nil{
		return nil,err
	}
	return workouts,nil



	
}




func (m WorkoutModel) GetWorkoutById(id int)([]Exercise,error){
	
	stmt:= `SELECT id,name,sets,reps
	FROM workout_exercises
	WHERE template_id=$1`

	ctx,cancel:= context.WithTimeout(context.Background(),3*time.Second)

	defer cancel()

	rows,err:= m.DB.QueryContext(ctx,stmt,id)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	var exercises []Exercise
	for rows.Next(){
		var e Exercise
		err= rows.Scan(&e.ID,&e.Name,&e.Sets,&e.Reps)
		if err!=nil{
			return nil,err
		}
		exercises= append(exercises,e)
	}
	if err= rows.Err(); err!=nil{
		return nil,err
	}
	return exercises,nil

}

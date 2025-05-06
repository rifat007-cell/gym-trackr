package data

import (
	"context"
	"database/sql"
	"time"
)

type Meal struct {
	ID int `json:"id"`
	Goal string `json:"goal"`
	DietaryPreference string `json:"dietary_preference"`
	Name string `json:"name"`
	Description string `json:"description"`
	Calories    int `json:"calories"`
}


type MealsModel struct {
	DB *sql.DB
}

func (m *MealsModel) GetAllMealByWorkoutName(goal, dietaryPreference string) ([]Meal, error) {
	query := `
		SELECT id, goal, dietary_preference, name, description, calories
		FROM meal_templates
		WHERE goal = $1 AND dietary_preference = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx,query, goal, dietaryPreference)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []Meal
	for rows.Next() {
		var meal Meal
		if err := rows.Scan(&meal.ID, &meal.Goal, &meal.DietaryPreference, &meal.Name, &meal.Description, &meal.Calories); err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return meals, nil
}

func (m *MealsModel) GetMealById(id int) (Meal, error) {
	query := `
		SELECT id, goal, dietary_preference, name, description, calories
		FROM meal_templates
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)

	var meal Meal
	if err := row.Scan(&meal.ID, &meal.Goal, &meal.DietaryPreference, &meal.Name, &meal.Description, &meal.Calories); err != nil {
		if err == sql.ErrNoRows {
			return meal, nil
		}
		return meal, err
	}

	return meal, nil
}
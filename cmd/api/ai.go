package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tanvir-rifat007/gymBuddy/internal/agents"
	"github.com/tanvir-rifat007/gymBuddy/internal/data"
)

func classifyIntent(msg string) string {
	lower := strings.ToLower(msg)
	switch {
	case containsProfanity(lower):
		return "offensive"
	case strings.Contains(lower, "meal"), strings.Contains(lower, "eat"), strings.Contains(lower, "diet"):
		return "meal_advice"
	case strings.Contains(lower, "exercise"), strings.Contains(lower, "workout"), strings.Contains(lower, "swap"):
		return "exercise_advice"
	case strings.Contains(lower, "plan"), strings.Contains(lower, "review"):
		return "full_review"
	case strings.Contains(lower, "tired"), strings.Contains(lower, "sore"), strings.Contains(lower, "pain"):
		return "fatigue"
	case strings.Contains(lower, "bored"), strings.Contains(lower, "same"), strings.Contains(lower, "routine"):
		return "boredom"
	case strings.Contains(lower, "motivate"), strings.Contains(lower, "motivated"), strings.Contains(lower, "lazy"):
		return "motivation"
	case strings.Contains(lower, "progress"), strings.Contains(lower, "track"):
		return "progress"
	case strings.Contains(lower, "supplement"), strings.Contains(lower, "creatine"), strings.Contains(lower, "protein powder"):
		return "supplements"
	case strings.Contains(lower, "more workout"), strings.Contains(lower, "add workout"):
		return "add_workout"
	case strings.Contains(lower, "change meal"), strings.Contains(lower, "new meal"):
		return "change_meal"
	case strings.Contains(lower, "custom workout"), strings.Contains(lower, "custom meal"):
		return "custom_plan"
	case lower == "hi" || lower == "hello" || lower == "hey":
		return "greeting"
	default:
		return "general"
	}
}

func containsProfanity(msg string) bool {
	profanities := []string{"fuck", "shit", "asshole", "bitch", "dumb", "stupid", "hate you", "fuck you", "screw you"}
	for _, word := range profanities {
		if strings.Contains(msg, word) {
			return true
		}
	}
	return false
}

func formatExercises(exercises []data.Exercise) string {
	var sb strings.Builder
	for _, ex := range exercises {
		sb.WriteString(fmt.Sprintf("- 🏋️ %s: %d sets x %d reps\n", ex.Name, ex.Sets, ex.Reps))
	}
	return sb.String()
}


// func (app *application) chat(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Message string `json:"message"`
// 	}

// 	if err := app.readJSON(w, r, &input); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	email, ok := r.Context().Value("email").(string)
// 	if !ok {
// 		app.serverErrorResponse(w, r, errors.New("missing email in context"))
// 		return
// 	}

// 	user, err := app.models.Users.GetUserByEmail(email)
// 	if err != nil {
// 		if errors.Is(err, data.ErrRecordNotFound) {
// 			app.notFoundResponse(w, r)
// 		} else {
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	userDailyPlan, err := app.models.UserDailyPlan.GetDailyPlanByUserID(user.ID)
// 	if err != nil {
// 		if errors.Is(err, data.ErrRecordNotFound) {
// 			app.notFoundResponse(w, r)
// 		} else {
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	exercises, err := app.models.Workouts.GetWorkoutById(userDailyPlan.WorkoutTemplateID)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	meal, err := app.models.Meals.GetMealById(userDailyPlan.MealTemplateID)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	intent := classifyIntent(input.Message)
// 	exerciseList := formatExercises(exercises)
// 	mealSummary := fmt.Sprintf("- 🍽️ %s: %s (%d cal)", meal.Name, meal.Description, meal.Calories)

// 	// Fetch last 10 chat messages
// 	history, err := app.models.ChatMessages.GetHistoryByUserID(user.ID, 10)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Prepare messages with context + history
// 	messages := []agents.OpenAPIMessage{
// 		{
// 			Role: "system",
// 			Content: fmt.Sprintf(`You are GymBuddy, a smart, supportive fitness coach AI.

// User's current plan:

// 🏋️ Workouts:
// %s

// 🍴 Meal Plan:
// %s

// Respond like a coach — concise, clear, and motivational. Adjust your advice to match user intent.`, exerciseList, mealSummary),
// 		},
// 	}

// 	for _, h := range history {
// 		messages = append(messages, agents.OpenAPIMessage{
// 			Role:    h.Role,
// 			Content: h.Content,
// 		})
// 	}

// 	// Append latest user message
// 	messages = append(messages, agents.OpenAPIMessage{
// 		Role:    "user",
// 		Content: fmt.Sprintf("Intent: %s\n\nMessage: %s", intent, input.Message),
// 	})

// 	// Run AI query
// 	response, err := app.ai.Query(messages)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	reply := response.Choices[0].Message.Content

// 	// Save messages after successful response
// 	_ = app.models.ChatMessages.Insert(user.ID, "user", input.Message)
// 	_ = app.models.ChatMessages.Insert(user.ID, "assistant", reply)

// 	// Send back AI reply
// 	app.writeJSON(w, http.StatusOK, envelope{"response": reply}, nil)
// }


func (app *application) chat(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Message string `json:"message"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("missing email in context"))
		return
	}

	user, err := app.models.Users.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	userDailyPlan, err := app.models.UserDailyPlan.GetDailyPlanByUserID(user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.writeJSON(w, http.StatusOK, envelope{
				"response": "👋 Before we chat, please visit the workout and meal pages to set up your daily plan. Once you have a workout and meal, I’ll be ready to help!",
			}, nil)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	exercises, err := app.models.Workouts.GetWorkoutById(userDailyPlan.WorkoutTemplateID)
	if err != nil {
		app.writeJSON(w, http.StatusOK, envelope{
			"response": "👋 It looks like you don’t have a workout assigned yet. Please select a workout from the workouts page, then come back and chat with me!",
		}, nil)
		return
	}

	meal, err := app.models.Meals.GetMealById(userDailyPlan.MealTemplateID)
	if err != nil {
		app.writeJSON(w, http.StatusOK, envelope{
			"response": "👋 It looks like you don’t have a meal plan yet. Please choose a meal from the meals page before we continue.",
		}, nil)
		return
	}

	intent := classifyIntent(input.Message)
	exerciseList := formatExercises(exercises)
	mealSummary := fmt.Sprintf("- 🍽️ %s: %s (%d cal)", meal.Name, meal.Description, meal.Calories)

	// Fetch last 10 chat messages
	history, err := app.models.ChatMessages.GetHistoryByUserID(user.ID, 10)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Prepare messages with context + history
	messages := []agents.OpenAPIMessage{
		{
			Role: "system",
			Content: fmt.Sprintf(`You are GymBuddy, a smart, supportive fitness coach AI.

User's current plan:

🏋️ Workouts:
%s

🍴 Meal Plan:
%s

Respond like a coach — concise, clear, and motivational. Adjust your advice to match user intent.`, exerciseList, mealSummary),
		},
	}

	for _, h := range history {
		messages = append(messages, agents.OpenAPIMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	// Append latest user message
	messages = append(messages, agents.OpenAPIMessage{
		Role:    "user",
		Content: fmt.Sprintf("Intent: %s\n\nMessage: %s", intent, input.Message),
	})

	// Run AI query
	response, err := app.ai.Query(messages)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	reply := response.Choices[0].Message.Content

	// Save messages after successful response
	_ = app.models.ChatMessages.Insert(user.ID, "user", input.Message)
	_ = app.models.ChatMessages.Insert(user.ID, "assistant", reply)

	// Send back AI reply
	app.writeJSON(w, http.StatusOK, envelope{"response": reply}, nil)
}




func (app *application) chatHistory(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("missing email in context"))
		return
	}

	user, err := app.models.Users.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	history, err := app.models.ChatMessages.GetHistoryByUserID(user.ID, 20) // last 20 messages
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"history": history}, nil)
}

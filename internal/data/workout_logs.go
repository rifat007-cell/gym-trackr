package data

import (
	"context"
	"database/sql"
	"time"
)


type WorkoutLog struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Exercise  string  `json:"exercise"`
	Sets       int     `json:"sets"`
	Reps      int     `json:"reps"`
	Duration   int     `json:"duration"`
	Weight     int `json:"weight"`
	LogDate   time.Time `json:"log_date"`
	CreatedAt time.Time `json:"created_at"`


}

type WorkoutLogModel struct {
	DB *sql.DB
}

type VolumeLog struct {
	LogDate     time.Time `json:"log_date"`
	TotalVolume int       `json:"total_volume"`
}


func (m *WorkoutLogModel) Insert(workoutLog *WorkoutLog) error {
	stmt:= `INSERT INTO user_workout_logs(user_id, workout_name, sets, reps, duration_minutes, weight_kg) VALUES($1, $2, $3, $4, $5, $6) RETURNING id,created_at,log_date`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{workoutLog.UserID, workoutLog.Exercise, workoutLog.Sets, workoutLog.Reps, workoutLog.Duration, workoutLog.Weight}

	err:= m.DB.QueryRowContext(ctx, stmt, args...).Scan(&workoutLog.ID, &workoutLog.CreatedAt, &workoutLog.LogDate)

	if err != nil {
		
		return err

	}

	return nil
}


func (m *WorkoutLogModel) GetVolumeOverTime(userID int) ([]VolumeLog, error) {
	stmt := `
		SELECT log_date, SUM(sets * reps * weight_kg) AS total_volume
		FROM user_workout_logs
		WHERE user_id = $1
		GROUP BY log_date
		ORDER BY log_date;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []VolumeLog
	for rows.Next() {
		var log VolumeLog
		if err := rows.Scan(&log.LogDate, &log.TotalVolume); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}



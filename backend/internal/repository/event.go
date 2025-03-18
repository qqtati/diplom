package repository

import (
	"backend/internal/models/event"
	"backend/pkg/storage"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type EventRepo struct {
	db storage.Postgres
}

func NewEventRepo(db storage.Postgres) *EventRepo {
	return &EventRepo{db: db}
}

// Проверяет, находится ли время события в рабочем времени учителя
func (r *EventRepo) isWithinWorkingHours(teacherID int64, startTime time.Time, duration int) (bool, error) {
	var workingHours string
	err := r.db.Get(&workingHours, `SELECT working_hours FROM "user" WHERE id = $1`, teacherID)
	if err != nil {
		return false, fmt.Errorf("error getting teacher working hours: %w", err)
	}

	if workingHours == "" {
		return false, fmt.Errorf("teacher has no working hours set")
	}

	// Получаем день недели события
	dayOfWeek := startTime.Weekday().String()
	// Получаем время начала и конца события
	eventStartTime := startTime.Add(3 * time.Hour).Format("15:04")
	eventEndTime := startTime.Add(3*time.Hour + time.Duration(duration)*time.Minute).Format("15:04")

	// Парсим JSON с рабочим временем
	var hours map[string]struct {
		Enabled bool   `json:"enabled"`
		Start   string `json:"start"`
		End     string `json:"end"`
	}
	if err := json.Unmarshal([]byte(workingHours), &hours); err != nil {
		return false, fmt.Errorf("error parsing working hours: %w", err)
	}

	// Проверяем, включен ли день недели и находится ли время в рабочем диапазоне
	dayHours, exists := hours[strings.ToLower(dayOfWeek)]
	if !exists || !dayHours.Enabled {
		return false, fmt.Errorf("teacher does not work on %s", dayOfWeek)
	}
	fmt.Println(eventStartTime, eventEndTime, dayHours.Start, dayHours.End)
	return eventStartTime >= dayHours.Start && eventEndTime <= dayHours.End, nil
}

func (r *EventRepo) GetTeacherEventsByUser(userId int64) ([]event.Event, error) {
	var result []event.Event
	err := r.db.Select(&result, `SELECT e.* FROM "event" e join "teacher_student" ts on e.teacher_id = ts.teacher_id where ts.student_id=$1`, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *EventRepo) GetEventsByUser(userId int64) ([]event.Event, error) {
	var result []event.Event
	err := r.db.Select(&result, `SELECT * FROM "event" e where e.teacher_id=$1 or e.student_id=$1`, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *EventRepo) GetEventById(id int64) (*event.Event, error) {
	var result []event.Event
	err := r.db.Select(&result, `SELECT * FROM "event" e where e.id=$1`, id)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}

func (r *EventRepo) InsertEvent(ev event.Event) (*int64, error) {
	// Проверяем рабочее время
	withinHours, err := r.isWithinWorkingHours(ev.TeacherID, ev.StartTime, ev.Duration)
	if err != nil {
		return nil, fmt.Errorf("error checking working hours: %w", err)
	}
	if !withinHours {
		return nil, fmt.Errorf("event time is outside teacher's working hours")
	}

	var id int64
	err = r.db.Get(&id, `INSERT INTO "event"(start_time, duration, teacher_id, price, student_id, description, is_recurring, recurrence_pattern, recurrence_end_date, parent_event_id, approved_by_teacher)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		ev.StartTime, ev.Duration, ev.TeacherID, ev.Price, ev.StudentID, ev.Description, ev.IsRecurring, ev.RecurrencePattern, ev.RecurrenceEndDate, ev.ParentEventID, ev.ApprovedByTeacher)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *EventRepo) DeleteEventById(eventId int64) error {
	_, err := r.db.Query(`DELETE FROM "event" WHERE id=$1 OR parent_event_id=$1`, eventId)
	return err
}

func (r *EventRepo) UpdateEvent(newEvent event.Event) error {
	// Проверяем рабочее время
	withinHours, err := r.isWithinWorkingHours(newEvent.TeacherID, newEvent.StartTime, newEvent.Duration)
	if err != nil {
		return fmt.Errorf("error checking working hours: %w", err)
	}
	if !withinHours {
		return fmt.Errorf("event time is outside teacher's working hours")
	}

	_, err = r.db.Query(`UPDATE "event" SET start_time=$2, duration=$3, price=$4, description=$5, student_id=$6, approved_by_teacher=$7, skipped=$8, rating=$9, is_recurring=$10, recurrence_pattern=$11, recurrence_end_date=$12 WHERE id=$1`,
		newEvent.ID, newEvent.StartTime, newEvent.Duration, newEvent.Price, newEvent.Description, newEvent.StudentID, newEvent.ApprovedByTeacher, newEvent.Skipped, newEvent.Rating, newEvent.IsRecurring, newEvent.RecurrencePattern, newEvent.RecurrenceEndDate)
	return err
}

func (r *EventRepo) GetRecurringEvents(parentEventId int64) ([]event.Event, error) {
	var result []event.Event
	err := r.db.Select(&result, `SELECT * FROM "event" WHERE parent_event_id=$1`, parentEventId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *EventRepo) GenerateRecurringEvents(parentEvent event.Event) error {
	if parentEvent.RecurrencePattern == "weekly" {
		currentDate := parentEvent.StartTime
		for currentDate.Before(parentEvent.RecurrenceEndDate) {
			currentDate = currentDate.Add(7 * 24 * time.Hour)

			// Проверяем рабочее время для каждого повторяющегося события
			withinHours, err := r.isWithinWorkingHours(parentEvent.TeacherID, currentDate, parentEvent.Duration)
			if err != nil {
				fmt.Printf("Error checking working hours for recurring event: %v\n", err)
				continue
			}
			if !withinHours {
				fmt.Printf("Skipping recurring event at %v - outside working hours\n", currentDate)
				continue
			}

			recurringEvent := event.Event{
				StartTime:         currentDate,
				Duration:          parentEvent.Duration,
				TeacherID:         parentEvent.TeacherID,
				Price:             parentEvent.Price,
				StudentID:         parentEvent.StudentID,
				Description:       parentEvent.Description,
				ApprovedByTeacher: parentEvent.ApprovedByTeacher,
				IsRecurring:       true,
				RecurrencePattern: parentEvent.RecurrencePattern,
				RecurrenceEndDate: parentEvent.RecurrenceEndDate,
				ParentEventID:     &parentEvent.ID,
			}
			_, err = r.InsertEvent(recurringEvent)
			if err != nil {
				fmt.Printf("Error creating recurring event: %v\n", err)
				continue
			}
		}
	}
	return nil
}

package event

import "time"

type Event struct {
	ID                int64     `json:"id" db:"id"`
	StartTime         time.Time `json:"start_time" db:"start_time"`
	Duration          int       `json:"duration" db:"duration"` // minutes
	TeacherID         int64     `json:"teacher_id" db:"teacher_id"`
	Price             int64     `json:"price" db:"price"` // per hour, rub
	StudentID         int64     `json:"student_id" db:"student_id"`
	Description       string    `json:"description" db:"description"`
	ApprovedByTeacher bool      `json:"approved_by_teacher" db:"approved_by_teacher"`
	Skipped           bool      `json:"skipped" db:"skipped"`
	Rating            *int      `json:"rating" db:"rating"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	IsRecurring       bool      `json:"is_recurring" db:"is_recurring"`
	RecurrencePattern string    `json:"recurrence_pattern" db:"recurrence_pattern"`
	RecurrenceEndDate time.Time `json:"recurrence_end_date" db:"recurrence_end_date"`
	ParentEventID     *int64    `json:"parent_event_id" db:"parent_event_id"`
}

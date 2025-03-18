package event

import "time"

type EventInput struct {
	Id          int64     `json:"id,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	TeacherID   int64     `json:"teacher_id"`
	Price       int64     `json:"price"`
	StudentID   int64     `json:"student_id,omitempty"`
	Description string    `json:"description,omitempty"`
	Skipped     bool      `json:"skipped,omitempty"`
	Rating      int       `json:"rating,omitempty"`
}

type DeleteEventInput struct {
	Id int64 `json:"id"`
}

type GetEventOutput struct {
	Id                int64     `json:"id"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	TeacherID         int64     `json:"teacher_id"`
	Price             int64     `json:"price"`
	StudentID         int64     `json:"student_id,omitempty"`
	Description       string    `json:"description,omitempty"`
	ApprovedByTeacher bool      `json:"approved_by_teacher,omitempty"`
	Skipped           bool      `json:"skipped,omitempty"`
	Rating            int       `json:"rating,omitempty"`
}

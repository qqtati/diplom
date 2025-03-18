package repository

import (
	"helprepet/internal/models/event"
	"helprepet/pkg/storage"
)

type EventRepo struct {
	db storage.Postgres
}

func NewEventRepo(db storage.Postgres) *EventRepo {
	return &EventRepo{db: db}
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
	var id int64
	err := r.db.Get(&id, `INSERT INTO "event"(start_time, duration, teacher_id, price, student_id, description)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		ev.StartTime, ev.Duration, ev.TeacherID, ev.Price, ev.StudentID, ev.Description)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *EventRepo) DeleteEventById(eventId int64) error {
	_, err := r.db.Query(`DELETE FROM "event" WHERE id=$1`, eventId)
	return err
}

func (r *EventRepo) UpdateEvent(newEvent event.Event) error {
	_, err := r.db.Query(`UPDATE "event" SET start_time=$2, duration=$3, price=$4, description=$5, student_id=$6, approved_by_teacher=$7, skipped=$8, rating=$9 WHERE id=$1`, newEvent.ID, newEvent.StartTime, newEvent.Duration, newEvent.Price, newEvent.Description, newEvent.StudentID, newEvent.ApprovedByTeacher, newEvent.Skipped, newEvent.Rating)
	return err
}

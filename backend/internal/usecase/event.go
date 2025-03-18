package usecase

import (
	"backend/internal/models/event"
	"backend/internal/repository"
	"fmt"
	"math"
	"time"
)

type EventUC struct {
	repo     *repository.EventRepo
	userRepo *repository.UserRepo
}

func NewEventUC(repo *repository.EventRepo, userRepo *repository.UserRepo) *EventUC {
	return &EventUC{
		userRepo: userRepo,
		repo:     repo,
	}
}

func (u *EventUC) GetEventsByUsername(username string) ([]event.GetEventOutput, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	events, err := u.repo.GetEventsByUser(user.ID)
	if err != nil {
		return nil, err
	}
	result := make([]event.GetEventOutput, len(events))
	for idx, ev := range events {
		result[idx] = event.GetEventOutput{
			Id:                ev.ID,
			StartTime:         ev.StartTime,
			EndTime:           ev.StartTime.Add(time.Minute * time.Duration(ev.Duration)),
			TeacherID:         ev.TeacherID,
			Price:             ev.Price,
			StudentID:         ev.StudentID,
			Description:       ev.Description,
			ApprovedByTeacher: ev.ApprovedByTeacher,
			Skipped:           ev.Skipped,
			Rating:            ev.Rating,
			IsRecurring:       ev.IsRecurring,
			RecurrencePattern: ev.RecurrencePattern,
			RecurrenceEndDate: ev.RecurrenceEndDate,
			ParentEventID:     ev.ParentEventID,
		}
	}
	return result, nil
}

func (u *EventUC) GetTeacherEvents(username string) ([]event.GetEventOutput, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	events, err := u.repo.GetTeacherEventsByUser(user.ID)
	if err != nil {
		return nil, err
	}
	result := make([]event.GetEventOutput, len(events))
	for idx, ev := range events {
		if ev.StudentID == user.ID {
			continue
		}
		result[idx] = event.GetEventOutput{
			Id:                ev.ID,
			StartTime:         ev.StartTime,
			EndTime:           ev.StartTime.Add(time.Minute * time.Duration(ev.Duration)),
			TeacherID:         ev.TeacherID,
			StudentID:         ev.StudentID,
			ApprovedByTeacher: ev.ApprovedByTeacher,
			Skipped:           ev.Skipped,
			Rating:            ev.Rating,
			IsRecurring:       ev.IsRecurring,
			RecurrencePattern: ev.RecurrencePattern,
			RecurrenceEndDate: ev.RecurrenceEndDate,
			ParentEventID:     ev.ParentEventID,
		}
	}
	return result, nil
}

func (u *EventUC) InsertEvent(ev event.EventInput) (*event.Event, error) {
	event := event.Event{
		StartTime:         ev.StartTime,
		Duration:          int(math.Ceil(ev.EndTime.Sub(ev.StartTime).Minutes())),
		TeacherID:         ev.TeacherID,
		Price:             ev.Price,
		StudentID:         ev.StudentID,
		Description:       ev.Description,
		Rating:            ev.Rating,
		IsRecurring:       ev.IsRecurring,
		RecurrencePattern: ev.RecurrencePattern,
		RecurrenceEndDate: ev.RecurrenceEndDate,
		ApprovedByTeacher: ev.ApprovedByTeacher,
	}

	id, err := u.repo.InsertEvent(event)
	if err != nil {
		return nil, err
	}

	// Если это повторяющееся событие, генерируем все его экземпляры
	if event.IsRecurring {
		event.ID = *id
		err = u.repo.GenerateRecurringEvents(event)
		if err != nil {
			return nil, err
		}
	}

	return u.repo.GetEventById(*id)
}

func (u *EventUC) DeleteEvent(ev event.DeleteEventInput) error {
	return u.repo.DeleteEventById(ev.Id)
}

func (u *EventUC) UpdateEvent(ev event.EventInput, username string) (*event.GetEventOutput, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	approvedByTeacher := false
	if user.ID == ev.TeacherID {
		approvedByTeacher = true
	}

	eve := event.Event{
		ID:                ev.Id,
		StartTime:         ev.StartTime,
		Duration:          int(math.Ceil(ev.EndTime.Sub(ev.StartTime).Minutes())),
		TeacherID:         ev.TeacherID,
		Price:             ev.Price,
		StudentID:         ev.StudentID,
		Description:       ev.Description,
		ApprovedByTeacher: approvedByTeacher,
		Skipped:           ev.Skipped,
		Rating:            ev.Rating,
		IsRecurring:       ev.IsRecurring,
		RecurrencePattern: ev.RecurrencePattern,
		RecurrenceEndDate: ev.RecurrenceEndDate,
	}

	err = u.repo.UpdateEvent(eve)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data, err := u.repo.GetEventById(ev.Id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	output := event.GetEventOutput{
		Id:                data.ID,
		StartTime:         data.StartTime,
		EndTime:           data.StartTime.Add(time.Minute * time.Duration(data.Duration)),
		TeacherID:         data.TeacherID,
		Price:             data.Price,
		StudentID:         data.StudentID,
		Description:       data.Description,
		Skipped:           data.Skipped,
		ApprovedByTeacher: data.ApprovedByTeacher,
		Rating:            data.Rating,
		IsRecurring:       data.IsRecurring,
		RecurrencePattern: data.RecurrencePattern,
		RecurrenceEndDate: data.RecurrenceEndDate,
		ParentEventID:     data.ParentEventID,
	}
	return &output, nil
}

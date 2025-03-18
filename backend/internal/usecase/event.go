package usecase

import (
	"fmt"
	"helprepet/internal/models/event"
	"helprepet/internal/repository"
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
		}
	}
	return result, nil
}

func (u *EventUC) InsertEvent(ev event.EventInput) (*event.Event, error) {
	id, err := u.repo.InsertEvent(event.Event{
		StartTime:   ev.StartTime,
		Duration:    int(math.Ceil(ev.EndTime.Sub(ev.StartTime).Minutes())),
		TeacherID:   ev.TeacherID,
		Price:       ev.Price,
		StudentID:   ev.StudentID,
		Description: ev.Description,
		Rating:      ev.Rating,
	})
	if err != nil {
		return nil, err
	}
	return u.repo.GetEventById(*id)
}

func (u *EventUC) DeleteEvent(ev event.DeleteEventInput) error {
	err := u.repo.DeleteEventById(ev.Id)
	return err
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
	err = u.repo.UpdateEvent(event.Event{
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
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	data, err := u.repo.GetEventById(ev.Id)
	fmt.Println(data, err)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return &event.GetEventOutput{
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
	}, nil
}

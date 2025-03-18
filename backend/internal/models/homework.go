package models

import (
	"time"
)

type Homework struct {
	ID          uint            `json:"id" db:"id"`
	Description string          `json:"description" db:"description"`
	DueDate     time.Time       `json:"dueDate" db:"due_date"`
	StudentID   uint            `json:"studentId" db:"student_id"`
	TeacherID   uint            `json:"teacherId" db:"teacher_id"`
	Rating      *int            `json:"rating" db:"rating"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`
	StudentName string          `json:"studentName" db:"student_name"`
	TeacherName string          `json:"teacherName" db:"teacher_name"`
	Files       []*HomeworkFile `json:"files" db:"-"`
}

type HomeworkFile struct {
	ID         uint      `json:"id" db:"id"`
	HomeworkID uint      `json:"homeworkId" db:"homework_id"`
	FileName   string    `json:"fileName" db:"file_name"`
	FilePath   string    `json:"filePath" db:"file_path"`
	UploadedAt time.Time `json:"uploadedAt" db:"uploaded_at"`
}

package service

import (
	"backend/internal/models"
	"backend/pkg/storage"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type HomeworkService struct {
	db storage.Postgres
}

func NewHomeworkService(db storage.Postgres) *HomeworkService {
	return &HomeworkService{db: db}
}

func (s *HomeworkService) CreateHomework(homework *models.Homework) error {
	fmt.Printf("Creating homework: %+v\n", homework)
	_, err := s.db.Exec(`
		INSERT INTO homeworks (description, due_date, student_id, teacher_id)
		VALUES ($1, $2, $3, $4)
	`, homework.Description, homework.DueDate, homework.StudentID, homework.TeacherID)
	if err != nil {
		fmt.Printf("Error creating homework: %v\n", err)
	}
	return err
}

func (s *HomeworkService) GetTeacherHomeworks(teacherID uint) ([]*models.Homework, error) {
	var homeworks []*models.Homework
	err := s.db.Select(&homeworks, `
		SELECT 
			h.id, 
			h.description, 
			h.due_date, 
			h.student_id, 
			h.teacher_id, 
			h.rating, 
			h.created_at, 
			h.updated_at, 
			u.name as student_name
		FROM homeworks h
		JOIN "user" u ON h.student_id = u.id
		WHERE h.teacher_id = $1
		ORDER BY h.due_date DESC
	`, teacherID)
	if err != nil {
		return nil, err
	}

	// Получаем файлы для каждого домашнего задания
	for _, homework := range homeworks {
		var files []*models.HomeworkFile
		err = s.db.Select(&files, `
			SELECT id, homework_id, file_name, file_path, uploaded_at
			FROM homework_files
			WHERE homework_id = $1
			ORDER BY uploaded_at DESC
		`, homework.ID)
		if err != nil {
			return nil, err
		}
		homework.Files = files
	}

	return homeworks, nil
}

func (s *HomeworkService) GetStudentHomeworks(studentID uint) ([]*models.Homework, error) {
	var homeworks []*models.Homework
	err := s.db.Select(&homeworks, `
		SELECT 
			h.id, 
			h.description, 
			h.due_date, 
			h.student_id, 
			h.teacher_id, 
			h.rating, 
			h.created_at, 
			h.updated_at, 
			u.name as teacher_name
		FROM homeworks h
		JOIN "user" u ON h.teacher_id = u.id
		WHERE h.student_id = $1
		ORDER BY h.due_date DESC
	`, studentID)
	if err != nil {
		return nil, err
	}

	// Получаем файлы для каждого домашнего задания
	for _, homework := range homeworks {
		var files []*models.HomeworkFile
		err = s.db.Select(&files, `
			SELECT id, homework_id, file_name, file_path, uploaded_at
			FROM homework_files
			WHERE homework_id = $1
			ORDER BY uploaded_at DESC
		`, homework.ID)
		if err != nil {
			return nil, err
		}
		homework.Files = files
	}

	return homeworks, nil
}

func (s *HomeworkService) UpdateRating(homeworkID uint, rating int) error {
	_, err := s.db.Exec(`
		UPDATE homeworks
		SET rating = $1
		WHERE id = $2
	`, rating, homeworkID)
	return err
}

func (s *HomeworkService) UploadFile(homeworkID uint, file *multipart.FileHeader) error {
	// Создаем директорию для файлов, если она не существует
	uploadDir := "uploads/homework"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию: %v", err)
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%d_%s%s", homeworkID, file.Filename, ext)
	filePath := filepath.Join(uploadDir, newFilename)

	// Открываем исходный файл
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer src.Close()

	// Создаем новый файл
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer dst.Close()

	// Копируем содержимое
	if _, err := dst.ReadFrom(src); err != nil {
		return fmt.Errorf("не удалось скопировать файл: %v", err)
	}

	// Сохраняем информацию о файле в БД
	_, err = s.db.Exec(`
		INSERT INTO homework_files (homework_id, file_name, file_path)
		VALUES ($1, $2, $3)
	`, homeworkID, file.Filename, filePath)
	return err
}

func (s *HomeworkService) GetFiles(homeworkID uint) ([]*models.HomeworkFile, error) {
	var files []*models.HomeworkFile
	err := s.db.Select(&files, `
		SELECT id, homework_id, file_name, file_path, created_at
		FROM homework_files
		WHERE homework_id = $1
		ORDER BY created_at DESC
	`, homeworkID)
	return files, err
}

func (s *HomeworkService) GetFileByID(fileID uint) (*models.HomeworkFile, error) {
	var file models.HomeworkFile
	err := s.db.Get(&file, `
		SELECT id, homework_id, file_name, file_path, uploaded_at
		FROM homework_files
		WHERE id = $1
	`, fileID)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти файл: %v", err)
	}

	// Проверяем существование файла
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл не найден на диске: %v", err)
	}

	return &file, nil
}

func (s *HomeworkService) GetFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

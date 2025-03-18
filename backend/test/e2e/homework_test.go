package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Homework struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	StudentID   int       `json:"student_id"`
	TeacherID   int       `json:"teacher_id"`
	Rating      *int      `json:"rating"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type HomeworkFile struct {
	ID         int       `json:"id"`
	HomeworkID int       `json:"homeworkId"`
	FileName   string    `json:"fileName"`
	FilePath   string    `json:"filePath"`
	UploadedAt time.Time `json:"uploadedAt"`
}

func TestHomeworkCRUD(t *testing.T) {
	timestamp := time.Now().Unix()

	// Создаем тестового учителя
	teacherUsername := fmt.Sprintf("teacher_%d", timestamp)
	registerUser(t, teacherUsername, "password123", "Test Teacher", 0)
	teacherToken := loginUser(t, teacherUsername, "password123")

	// Создаем тестового ученика
	studentUsername := fmt.Sprintf("student_%d", timestamp)
	registerUser(t, studentUsername, "password123", "Test Student", 1)
	studentToken := loginUser(t, studentUsername, "password123")

	// Создаем связь учитель-ученик
	createTeacherStudentRelation(t, teacherToken, studentToken)

	// Тест создания домашнего задания
	homework := createHomework(t, teacherToken, studentToken)

	// Тест получения домашних заданий учителем
	getTeacherHomeworks(t, teacherToken)

	// Тест получения домашних заданий учеником
	getStudentHomeworks(t, studentToken)

	// Тест загрузки файла
	file := uploadHomeworkFile(t, studentToken, homework.ID)

	// Тест получения файлов
	getHomeworkFiles(t, studentToken, homework.ID)

	// Тест обновления оценки
	updateHomeworkRating(t, teacherToken, homework.ID)

	// Удаляем тестовые файлы
	os.Remove(file.FilePath)
}

func createTeacherStudentRelation(t *testing.T, teacherToken, studentToken string) {
	// Получаем ID ученика
	studentID := getUserID(t, studentToken)

	// Создаем связь
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/teacher_student", BaseURL), nil)
	assert.NoError(t, err)

	req.Header.Set("AccessToken", teacherToken)
	q := req.URL.Query()
	q.Add("student_id", fmt.Sprintf("%d", studentID))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func getUserID(t *testing.T, token string) int {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/me", BaseURL), nil)
	assert.NoError(t, err)

	req.Header.Set("AccessToken", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Status string `json:"status"`
		Result struct {
			ID int64 `json:"id"`
		} `json:"result"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "SUCCESS", result.Status)

	return int(result.Result.ID)
}

func createHomework(t *testing.T, teacherToken string, studentToken string) *Homework {
	// Получаем ID ученика
	studentID := getUserID(t, studentToken)

	homeworkData := map[string]interface{}{
		"description": "Тестовое домашнее задание",
		"due_date":    time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		"student_id":  studentID,
	}

	jsonData, err := json.Marshal(homeworkData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/homework/", BaseURL), bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AccessToken", teacherToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("Homework POST response status: %d, body: %s", resp.StatusCode, string(body))

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var homework Homework
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&homework)
	assert.NoError(t, err)

	return &homework
}

func getTeacherHomeworks(t *testing.T, teacherToken string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/homework/", BaseURL), nil)
	assert.NoError(t, err)

	req.Header.Set("AccessToken", teacherToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var homeworks []Homework
	err = json.NewDecoder(resp.Body).Decode(&homeworks)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeworks)
}

func getStudentHomeworks(t *testing.T, studentToken string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/homework/", BaseURL), nil)
	assert.NoError(t, err)

	req.Header.Set("AccessToken", studentToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var homeworks []Homework
	err = json.NewDecoder(resp.Body).Decode(&homeworks)
	assert.NoError(t, err)
	assert.NotEmpty(t, homeworks)
}

func uploadHomeworkFile(t *testing.T, studentToken string, homeworkID int) *HomeworkFile {
	// Создаем временный файл
	file, err := os.CreateTemp("", "test_homework_*.txt")
	assert.NoError(t, err)
	defer file.Close()

	// Записываем тестовые данные в файл
	_, err = file.WriteString("Тестовое содержимое файла")
	assert.NoError(t, err)

	// Создаем multipart запрос
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Добавляем файл
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	// Закрываем writer
	err = writer.Close()
	assert.NoError(t, err)

	// Создаем запрос
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/homework/%d/file", BaseURL, homeworkID), &buf)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("AccessToken", studentToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var homeworkFile HomeworkFile
	err = json.NewDecoder(resp.Body).Decode(&homeworkFile)
	assert.NoError(t, err)

	return &homeworkFile
}

func getHomeworkFiles(t *testing.T, studentToken string, homeworkID int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/homework/%d/file", BaseURL, homeworkID), nil)
	assert.NoError(t, err)

	req.Header.Set("AccessToken", studentToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var files []HomeworkFile
	err = json.NewDecoder(resp.Body).Decode(&files)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)
}

func updateHomeworkRating(t *testing.T, teacherToken string, homeworkID int) {
	ratingData := map[string]interface{}{
		"rating": 5,
	}

	jsonData, err := json.Marshal(ratingData)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/homework/%d/rating", BaseURL, homeworkID), bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AccessToken", teacherToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

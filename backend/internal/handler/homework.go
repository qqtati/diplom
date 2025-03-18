package handler

import (
	"backend/internal/common"
	"backend/internal/models"
	"backend/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HomeworkHandler struct {
	homeworkService *service.HomeworkService
}

func NewHomeworkHandler(homeworkService *service.HomeworkService) *HomeworkHandler {
	return &HomeworkHandler{
		homeworkService: homeworkService,
	}
}

func (h *HomeworkHandler) GetHomeworkService() *service.HomeworkService {
	return h.homeworkService
}

func (h *HomeworkHandler) CreateHomework(c *fiber.Ctx) error {
	var req struct {
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
		StudentID   uint   `json:"student_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	// Получаем ID учителя из контекста
	teacherID := c.Locals("userID")
	if teacherID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Не авторизован",
		})
	}

	// Преобразуем teacherID в uint
	var teacherIDUint uint
	switch v := teacherID.(type) {
	case int64:
		teacherIDUint = uint(v)
	case int:
		teacherIDUint = uint(v)
	case uint:
		teacherIDUint = v
	default:
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Неверный тип ID учителя",
		})
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат даты",
		})
	}

	homework := &models.Homework{
		Description: req.Description,
		DueDate:     dueDate,
		StudentID:   req.StudentID,
		TeacherID:   teacherIDUint,
	}

	if err := h.homeworkService.CreateHomework(homework); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось создать домашнее задание",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"result": homework,
	})
}

func (h *HomeworkHandler) GetHomeworks(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Не авторизован",
		})
	}

	userRole := c.Locals("userRole")
	if userRole == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Не авторизован",
		})
	}

	var homeworks []*models.Homework
	var err error

	if userRole == 0 { // Учитель
		homeworks, err = h.homeworkService.GetTeacherHomeworks(uint(userID.(int64)))
	} else { // Ученик
		homeworks, err = h.homeworkService.GetStudentHomeworks(uint(userID.(int64)))
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить домашние задания",
		})
	}

	return c.JSON(common.Response{
		Result: homeworks,
	})
}

func (h *HomeworkHandler) UpdateRating(c *fiber.Ctx) error {
	var req struct {
		HomeworkID uint `json:"homework_id"`
		Rating     int  `json:"rating"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	if err := h.homeworkService.UpdateRating(req.HomeworkID, req.Rating); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить оценку",
		})
	}

	return c.JSON(fiber.Map{
		"result": "Оценка обновлена",
	})
}

func (h *HomeworkHandler) UploadFile(c *fiber.Ctx) error {
	homeworkID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID задания",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Файл не найден",
		})
	}

	if err := h.homeworkService.UploadFile(uint(homeworkID), file); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось загрузить файл",
		})
	}

	return c.JSON(fiber.Map{
		"result": "Файл загружен",
	})
}

func (h *HomeworkHandler) GetFiles(c *fiber.Ctx) error {
	homeworkID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID задания",
		})
	}

	files, err := h.homeworkService.GetFiles(uint(homeworkID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить файлы",
		})
	}

	return c.JSON(fiber.Map{
		"result": files,
	})
}

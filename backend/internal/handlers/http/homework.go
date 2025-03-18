package http

import (
	"backend/internal/common"
	"backend/internal/handler"
	"backend/internal/models"
	"backend/internal/service"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HomeworkHandler struct {
	handler *handler.HomeworkHandler
}

func NewHomeworkHandler(homeworkService *service.HomeworkService) *HomeworkHandler {
	return &HomeworkHandler{
		handler: handler.NewHomeworkHandler(homeworkService),
	}
}

func (h *HomeworkHandler) getHomeworkService() *service.HomeworkService {
	return h.handler.GetHomeworkService()
}

func MapHomeworkRoutes(router fiber.Router, h *HomeworkHandler, userHandler *UserHandler) {
	router.Post("/", userHandler.CheckAccessToken(), h.CreateHomework)
	router.Get("/", userHandler.CheckAccessToken(), h.GetHomeworks)
	router.Put("/:id/rating", userHandler.CheckAccessToken(), h.UpdateRating)
	router.Post("/:id/file", userHandler.CheckAccessToken(), h.UploadFile)
	router.Get("/:id/file", userHandler.CheckAccessToken(), h.GetFiles)
	router.Get("/file/:id", userHandler.CheckAccessToken(), h.DownloadFile)
}

func (h *HomeworkHandler) CreateHomework(c *fiber.Ctx) error {
	var req struct {
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
		StudentID   uint   `json:"student_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	// Получаем ID учителя из контекста
	teacherID := c.Locals("userID")
	if teacherID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Неверный тип ID учителя",
		})
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат даты",
		})
	}

	homework := &models.Homework{
		Description: req.Description,
		DueDate:     dueDate,
		StudentID:   req.StudentID,
		TeacherID:   teacherIDUint,
	}

	if err := h.getHomeworkService().CreateHomework(homework); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось создать домашнее задание",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(homework)
}

func (h *HomeworkHandler) GetHomeworks(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Не авторизован",
		})
	}

	userRole := c.Locals("userRole")
	if userRole == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Не авторизован",
		})
	}

	var homeworks []*models.Homework
	var err error
	fmt.Println(userRole)
	if userRole.(int64) == 0 { // Учитель
		homeworks, err = h.getHomeworkService().GetTeacherHomeworks(uint(userID.(int64)))
	} else { // Ученик
		homeworks, err = h.getHomeworkService().GetStudentHomeworks(uint(userID.(int64)))
	}
	fmt.Println(userID, userRole, homeworks, err)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить домашние задания",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.Response{
		Result: homeworks,
	})
}

func (h *HomeworkHandler) UpdateRating(c *fiber.Ctx) error {
	var req struct {
		HomeworkID uint `json:"homework_id"`
		Rating     int  `json:"rating"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	if err := h.getHomeworkService().UpdateRating(req.HomeworkID, req.Rating); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID задания",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Файл не найден",
		})
	}

	if err := h.getHomeworkService().UploadFile(uint(homeworkID), file); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID задания",
		})
	}

	files, err := h.getHomeworkService().GetFiles(uint(homeworkID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить файлы",
		})
	}

	return c.Status(fiber.StatusOK).JSON(files)
}

func (h *HomeworkHandler) DownloadFile(c *fiber.Ctx) error {
	fileID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID файла",
		})
	}

	file, err := h.getHomeworkService().GetFileByID(uint(fileID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	content, err := h.getHomeworkService().GetFileContent(file.FilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось прочитать файл",
		})
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.FileName))
	return c.Send(content)
}

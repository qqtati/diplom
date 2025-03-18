package http

import (
	"backend/internal/common"
	"backend/internal/models/user"
	"backend/internal/usecase"
	util_http "backend/pkg/util/http"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	uc      *usecase.UserUC
	reqUtil *util_http.Reader
}

func NewUserHandler(uc *usecase.UserUC, reqUtil *util_http.Reader) *UserHandler {
	return &UserHandler{
		uc:      uc,
		reqUtil: reqUtil,
	}
}

func (h UserHandler) SignIn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body user.AuthorizeInput
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return common.ErrInput.Wrap(err)
		}

		result, err := h.uc.AuthorizeUser(&body)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status: common.FailedStatus,
				Result: nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: *result,
		})
	}
}

func (h UserHandler) SignUp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body user.SignUpInput
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &body)
		if err != nil {
			return common.ErrInput.Wrap(err)
		}

		result, err := h.uc.SignUp(&body)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: *result,
		})
	}
}

func (h UserHandler) GetMe() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.UserContext().Value("username").(string)
		result, err := h.uc.GetMe(username)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: *result,
		})
	}
}

func (h UserHandler) GetStudents() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.UserContext().Value("username").(string)
		result, err := h.uc.GetStudents(username)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h UserHandler) GetStudentsStats() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.UserContext().Value("username").(string)
		days := ctx.QueryInt("days", -1)
		result, err := h.uc.GetStudentStats(username, days)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h UserHandler) GetTeachers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.UserContext().Value("username").(string)
		result, err := h.uc.GetTeachers(username)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h UserHandler) GetTeachersStats() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.UserContext().Value("username").(string)
		days := ctx.QueryInt("days", -1)
		result, err := h.uc.GetTeacherStats(username, days)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: result,
		})
	}
}

func (h UserHandler) CheckAccessToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var token = ctx.Get("AccessToken")
		username, err := h.uc.CheckToken(token)
		if err != nil {
			return common.ErrUnauthorized.Wrap(err)
		}
		if username == nil {
			return common.ErrUnauthorized
		}

		// Получаем информацию о пользователе
		user, err := h.uc.GetMe(*username)
		if err != nil {
			return common.ErrUnauthorized.Wrap(err)
		}

		ctx.Locals("userID", user.ID)
		ctx.Locals("userRole", user.Role)
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), "username", *username))
		return ctx.Next()
	}
}

func (h UserHandler) CreateTeacherStudentRelation() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Получаем ID учителя из контекста
		teacherUsername := ctx.UserContext().Value("username").(string)
		teacher, err := h.uc.GetMe(teacherUsername)
		if err != nil {
			return common.ErrUnauthorized.Wrap(err)
		}

		// Проверяем, что пользователь является учителем
		if teacher.Role != 1 {
			return ctx.Status(fiber.StatusForbidden).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: "Только учителя могут добавлять учеников",
			})
		}

		// Получаем ID ученика из query параметра
		studentID := ctx.QueryInt("student_id", 0)
		if studentID == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: "ID ученика не указан",
			})
		}

		// Создаем связь
		err = h.uc.CreateTeacherStudentRelation(uint(teacher.ID), uint(studentID))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: "Не удалось создать связь учитель-ученик",
			})
		}

		return ctx.Status(fiber.StatusCreated).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func (h UserHandler) GetProfile() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.UserContext().Value("username").(string)
		result, err := h.uc.GetProfile(username)
		if err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
			Result: *result,
		})
	}
}

func (h UserHandler) UpdateProfile() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input user.UserInput
		if err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &input); err != nil {
			return common.ErrInput.Wrap(err)
		}

		username := ctx.UserContext().Value("username").(string)
		user, err := h.uc.GetMe(username)
		if err != nil {
			return common.ErrUnauthorized.Wrap(err)
		}

		if err := h.uc.UpdateProfile(ctx.Context(), user.ID, input); err != nil {
			var logicErr common.LogicError
			ok := errors.As(err, &logicErr)
			if !ok {
				return err
			}

			return ctx.Status(logicErr.Code).JSON(common.Response{
				Status:      common.FailedStatus,
				Description: logicErr.Message,
				Result:      nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(common.Response{
			Status: common.SuccessStatus,
		})
	}
}

func MapRoutes(r fiber.Router, h *UserHandler) {
	r.Post("/sign_in", h.SignIn())
	r.Post("/sign_up", h.SignUp())

	r.Get("/me", h.CheckAccessToken(), h.GetMe())
	r.Get("/students", h.CheckAccessToken(), h.GetStudents())
	r.Get("/students/stats", h.CheckAccessToken(), h.GetStudentsStats())

	r.Get("/teachers", h.CheckAccessToken(), h.GetTeachers())
	r.Get("/teachers/stats", h.CheckAccessToken(), h.GetTeachersStats())

	r.Post("/teacher_student", h.CheckAccessToken(), h.CreateTeacherStudentRelation())

	r.Get("/profile", h.CheckAccessToken(), h.GetProfile())
	r.Put("/profile", h.CheckAccessToken(), h.UpdateProfile())
}

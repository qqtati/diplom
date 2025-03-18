package http

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"helprepet/internal/common"
	"helprepet/internal/models/user"
	"helprepet/internal/usecase"
	util_http "helprepet/pkg/util/http"
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
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), "username", *username))
		return ctx.Next()
	}
}

func MapRoutes(r fiber.Router, h *UserHandler) {
	r.Post("/sign_in", h.SignIn())
	r.Post("/sign_up", h.SignUp())

	r.Get("/me", h.CheckAccessToken(), h.GetMe())
	r.Get("/students", h.CheckAccessToken(), h.GetStudents())
	r.Get("/students/stats", h.CheckAccessToken(), h.GetStudentsStats())
}

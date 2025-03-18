package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"helprepet/internal/common"
	"helprepet/internal/models/event"
	"helprepet/internal/usecase"
	util_http "helprepet/pkg/util/http"
)

type EventHandler struct {
	uc      *usecase.EventUC
	reqUtil *util_http.Reader
}

func NewEventHandler(uc *usecase.EventUC, reqUtil *util_http.Reader) *EventHandler {
	return &EventHandler{
		uc:      uc,
		reqUtil: reqUtil,
	}
}

func (h EventHandler) CreateEvent() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var data event.EventInput
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &data)
		if err != nil {
			return err
		}
		result, err := h.uc.InsertEvent(data)
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

func (h EventHandler) UpdateEvent() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var data event.EventInput
		username := ctx.UserContext().Value("username").(string)
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &data)
		if err != nil {
			return err
		}
		result, err := h.uc.UpdateEvent(data, username)
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

func (h EventHandler) DeleteEvent() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var data event.DeleteEventInput
		err := h.reqUtil.Read(ctx.Context(), ctx.BodyParser, &data)
		if err != nil {
			return err
		}
		err = h.uc.DeleteEvent(data)
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
		})
	}
}

func (h EventHandler) GetEvents() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.UserContext().Value("username").(string)
		result, err := h.uc.GetEventsByUsername(username)
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

func MapEventRoutes(r fiber.Router, h *EventHandler, mw *UserHandler) {
	r.Post("/", mw.CheckAccessToken(), h.CreateEvent())
	r.Put("/", mw.CheckAccessToken(), h.UpdateEvent())
	r.Delete("/", mw.CheckAccessToken(), h.DeleteEvent())
	r.Get("/", mw.CheckAccessToken(), h.GetEvents())
}

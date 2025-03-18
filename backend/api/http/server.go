package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"helprepet/api"
	api_http_model "helprepet/api/http/model"
	"helprepet/config"
	"helprepet/internal"
	"helprepet/internal/common"
	"helprepet/internal/handlers/http"
	"helprepet/internal/usecase"
	"helprepet/pkg/util/http"

	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"

	goJson "github.com/goccy/go-json"
)

type httpServer struct {
	fiber *fiber.App
	cfg   *config.Config
}

func NewHttpServer(cfg *config.Config) api.Server {
	return &httpServer{
		cfg: cfg,
	}
}

func (h *httpServer) Init() error {
	h.fiber = fiber.New(fiber.Config{
		Immutable:               true,
		AppName:                 "main",
		EnableTrustedProxyCheck: true,
		JSONEncoder:             goJson.Marshal,
		JSONDecoder:             goJson.Unmarshal,
	})

	h.fiber.Use(recoverMDW.New(recoverMDW.Config{
		EnableStackTrace: true,
	}))

	return nil
}

func (h *httpServer) MapHandlers(app *internal.App) error {
	// ENGINE
	h.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	h.fiber.Get("/version", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(api_http_model.VersionResponse{
			Version: h.cfg.Server.Version,
			Response: common.Response{
				Status: common.SuccessStatus,
			},
		})
	})

	// UTILS
	reqReader := util_http.NewReader()

	// HANDLERS
	userHandler := http.NewUserHandler(app.UC["user"].(*usecase.UserUC), reqReader)
	eventHandler := http.NewEventHandler(app.UC["event"].(*usecase.EventUC), reqReader)

	userGroup := h.fiber.Group("/user")
	http.MapRoutes(userGroup, userHandler)

	eventGroup := h.fiber.Group("/event")
	http.MapEventRoutes(eventGroup, eventHandler, userHandler)
	return nil
}

func (h *httpServer) Run() error {
	fmt.Printf("LISTENING %s:%s\n", h.cfg.Server.Host, h.cfg.Server.Port)
	err := h.fiber.Listen(h.cfg.Server.Host + ":" + h.cfg.Server.Port)
	return err
}

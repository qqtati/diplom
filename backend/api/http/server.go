package http

import (
	"backend/api"
	api_http_model "backend/api/http/model"
	"backend/config"
	"backend/internal"
	"backend/internal/common"
	"backend/internal/handlers"
	"backend/internal/handlers/http"
	"backend/internal/service"
	"backend/internal/usecase"
	util_http "backend/pkg/util/http"
	"fmt"

	"github.com/gofiber/fiber/v2"
	loggerMDW "github.com/gofiber/fiber/v2/middleware/logger"

	goJson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"
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
	h.fiber.Use(loggerMDW.New(loggerMDW.Config{}))

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
	wsHandler := http.NewWebSocketHandler()
	homeworkHandler := http.NewHomeworkHandler(app.UC["homework"].(*service.HomeworkService))

	userGroup := h.fiber.Group("/user")
	http.MapRoutes(userGroup, userHandler)

	eventGroup := h.fiber.Group("/event")
	http.MapEventRoutes(eventGroup, eventHandler, userHandler)

	// WebSocket routes
	wsGroup := h.fiber.Group("/ws")
	http.MapWebSocketRoutes(wsGroup, wsHandler, userHandler)

	// WebRTC routes
	handlers.SetupWebRTCRoutes(h.fiber)

	// Homework routes
	homeworkGroup := h.fiber.Group("/homework")
	http.MapHomeworkRoutes(homeworkGroup, homeworkHandler, userHandler)

	return nil
}

func (h *httpServer) Run() error {
	fmt.Printf("LISTENING %s:%s\n", h.cfg.Server.Host, h.cfg.Server.Port)
	err := h.fiber.Listen(h.cfg.Server.Host + ":" + h.cfg.Server.Port)
	return err
}

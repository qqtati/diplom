package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/handler"
	"backend/internal/models"
	"backend/internal/service"
)

func main() {
	// Подключение к базе данных
	dsn := "host=localhost user=postgres password=postgres dbname=tutor port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Автоматическая миграция моделей
	err = db.AutoMigrate(&models.Homework{}, &models.HomeworkFile{})
	if err != nil {
		log.Fatalf("Не удалось выполнить миграцию: %v", err)
	}

	// Инициализация сервисов
	homeworkService := service.NewHomeworkService(db)

	// Инициализация обработчиков
	homeworkHandler := handler.NewHomeworkHandler(homeworkService)

	// Настройка роутера
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, AccessToken")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Группа маршрутов для домашних заданий
	homework := r.Group("/homework")
	{
		homework.POST("/", homeworkHandler.CreateHomework)
		homework.GET("/", homeworkHandler.GetHomeworks)
		homework.PUT("/:id/rating", homeworkHandler.UpdateRating)
		homework.POST("/:id/file", homeworkHandler.UploadFile)
		homework.GET("/:id/files", homeworkHandler.GetFiles)
	}

	// Запуск сервера
	if err := r.Run(":9001"); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

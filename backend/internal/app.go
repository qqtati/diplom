package internal

import (
	"fmt"
	"github.com/go-redis/redis"
	"helprepet/config"
	"helprepet/internal/repository"
	"helprepet/internal/usecase"
	"helprepet/pkg/storage"
)

type App struct {
	cfg   *config.Config
	UC    map[string]interface{}
	Repo  map[string]interface{}
	db    storage.Postgres
	redis *storage.Redis
}

func NewApp(cfg *config.Config) *App {
	return &App{
		UC:   make(map[string]interface{}),
		Repo: make(map[string]interface{}),
		cfg:  cfg,
	}
}

func (a *App) Init() error {
	// Storage init
	var err error
	a.db, err = storage.Connect(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		a.cfg.Postgres.Host,
		a.cfg.Postgres.Port,
		a.cfg.Postgres.User,
		a.cfg.Postgres.Password,
		a.cfg.Postgres.DBName,
		a.cfg.Postgres.SSLMode))
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     a.cfg.Redis.Host + ":" + a.cfg.Redis.Port,
		Password: a.cfg.Redis.Password,
	})
	if err := redisClient.Ping().Err(); err != nil {
		return err
	}
	a.redis = storage.NewRedis(redisClient)
	// Repo init
	a.Repo["user"] = repository.NewUserRepo(a.db)
	a.Repo["event"] = repository.NewEventRepo(a.db)
	// UC init
	a.UC["user"] = usecase.NewUserUC(a.Repo["user"].(*repository.UserRepo), a.redis, &a.cfg.Auth, &a.cfg.Mail)
	a.UC["event"] = usecase.NewEventUC(a.Repo["event"].(*repository.EventRepo), a.Repo["user"].(*repository.UserRepo))
	return nil
}

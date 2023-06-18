package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/twirp_handlers"
	cfg "github.com/satont/tsuwari/libs/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"net/http"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	config, err := cfg.New()
	if err != nil {
		logger.Sugar().Panic(err)
	}

	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Panic("failed to connect database", err)
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		logger.Sugar().Panic(err)
	}
	redisClient := redis.NewClient(redisOpts)

	twirpProtectedPath, twirpProtectedHandler := twirp_handlers.NewProtected(twirp_handlers.Opts{
		Redis: redisClient,
		DB:    db,
	})

	mux := http.NewServeMux()
	mux.Handle(twirpProtectedPath, twirpProtectedHandler)
	logger.Sugar().Panic(http.ListenAndServe(":3002", mux))
}

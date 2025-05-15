package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"uala-timeline-service/internal/domain"
	"uala-timeline-service/internal/infrastructure"
	"uala-timeline-service/libs/events"
)

type Dependencies struct {
	EventPublisher  events.Publisher
	TimelineService domain.TimelineFilledService
}

func BuildDependencies(config Config) (*Dependencies, error) {
	natsPublisher := events.NewNatsPublisher()
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
	)
	if !config.Postgres.UseSSL {
		url += "?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Redis.Host,
	})
	res := redisClient.Ping(context.Background())
	if res.Err() != nil {
		panic(res.Err())
	}

	timelineRepository := infrastructure.NewTimelineRepository(db)
	timelineFilledRepository := infrastructure.NewTimelineFilledRepository(redisClient)

	timelineService := domain.NewTimelineService(timelineRepository, nil, timelineFilledRepository)

	return &Dependencies{
		TimelineService: timelineService,
		EventPublisher:  natsPublisher,
	}, nil
}

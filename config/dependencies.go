package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"uala-timeline-service/internal/domain"
	"uala-timeline-service/internal/infrastructure"
	"uala-timeline-service/libs/events"
)

type Dependencies struct {
	EventPublisher   events.Publisher
	FollowRepository domain.FollowRepository
	PostRepository   domain.PostRepository
	TimelineService  domain.DayUserTimelineFilledService
}

func BuildDependencies(config Config) (*Dependencies, error) {
	// Nats boot
	natsPublisher := events.NewNatsPublisher()

	// Postgres boot
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

	// Redis boot
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Redis.Host,
	})
	res := redisClient.Ping(context.Background())
	if res.Err() != nil {
		panic(res.Err())
	}

	// AWS boot
	var a = "http://localhost:8000"

	awsCfg := aws.Config{
		Region: config.AWS.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			config.AWS.Account,
			config.AWS.Secret,
			"",
		),
		BaseEndpoint: &a,
	}

	dynamoDb := dynamodb.NewFromConfig(awsCfg)

	timelineRepository := infrastructure.NewTimelineRepository(db)
	dayTimelineFilledRepository := infrastructure.NewDynamoPaymentRepository(dynamoDb, config.AWS.Table)
	postRepository := infrastructure.NewRestPostRepository("http://localhost:8080")
	followsRepository := infrastructure.NewRestFollowsRepository("http://localhost:8082")

	timelineService := domain.NewTimelineService(timelineRepository, postRepository, dayTimelineFilledRepository)

	return &Dependencies{
		TimelineService:  timelineService,
		EventPublisher:   natsPublisher,
		FollowRepository: followsRepository,
		PostRepository:   postRepository,
	}, nil
}

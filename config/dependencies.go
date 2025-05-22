package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"uala-timeline-service/internal/domain/day_timeline_filled/service"
	"uala-timeline-service/internal/domain/follows"
	"uala-timeline-service/internal/domain/posts"
	"uala-timeline-service/internal/infrastructure"
	"uala-timeline-service/libs/events"
)

type Dependencies struct {
	EventPublisher   events.Publisher
	FollowRepository follows.FollowRepository
	PostRepository   posts.PostRepository
	TimelineService  service.DayUserTimelineFilledService
}

func BuildDependencies(config Config) (*Dependencies, error) {
	// Nats boot
	natsPublisher := events.NewNatsPublisher(config.Nats.Host)

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
	postRepository := infrastructure.NewRestPostRepository(config.RestConfigs.PostService.BasePath)
	followsRepository := infrastructure.NewRestFollowsRepository(config.RestConfigs.FollowersService.BasePath)

	timelineService := service.NewTimelineService(timelineRepository, postRepository, dayTimelineFilledRepository)

	return &Dependencies{
		TimelineService:  timelineService,
		EventPublisher:   natsPublisher,
		FollowRepository: followsRepository,
		PostRepository:   postRepository,
	}, nil
}

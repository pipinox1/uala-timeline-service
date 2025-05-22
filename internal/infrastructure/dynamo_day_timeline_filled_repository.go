package infrastructure

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"io"
	"time"
	"uala-timeline-service/internal/domain/day_timeline_filled"
	"uala-timeline-service/internal/domain/posts"
)

var _ day_timeline_filled.DayUserTimelineFilledRepository = (*DynamoDayTimelineFilledRepository)(nil)

var emptyDayFilled = day_timeline_filled.DayUserTimelineFilled{}

const (
	dayPrefix = "timeline:"
	pkPrefix  = "user:%s"
	skPrefix  = "day:%s"
)

type DynamoDayTimelineFilledRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoPaymentRepository(client *dynamodb.Client, tableName string) *DynamoDayTimelineFilledRepository {
	return &DynamoDayTimelineFilledRepository{
		client:    client,
		tableName: tableName,
	}
}

func (d *DynamoDayTimelineFilledRepository) GetDayUserTimelineFilled(ctx context.Context, filter day_timeline_filled.DayUserTimelineFilledFilter) (*day_timeline_filled.DayUserTimelineFilled, error) {
	dayKey := fmt.Sprintf("%v:%v:%v", filter.FromYear, filter.FromMonth, filter.FromDay)
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: buildPK(filter.UserID)},
			"sk": &types.AttributeValueMemberS{Value: buildSK(dayKey)},
		},
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return &day_timeline_filled.DayUserTimelineFilled{
			Posts:  nil,
			UserID: filter.UserID,
		}, nil
	}

	var dayTimeline DynamoDayUserTimelinePage
	err = attributevalue.UnmarshalMap(result.Item, &dayTimeline)
	if err != nil {
		return nil, err
	}

	return dayTimeline.toDomain()
}

func (d *DynamoDayTimelineFilledRepository) AddPosts(ctx context.Context, userID string, post []posts.Post) error {
	dayPostMap := splitPostByDate(post)
	var transactItems []types.TransactWriteItem
	for dayKey, newPosts := range dayPostMap {
		oldDynamoDayFilled, err := d.getDayFilled(ctx, userID, dayKey)
		if err != nil {
			return err
		}

		newPostsCompressed := make([]string, len(newPosts))
		for i, post := range newPosts {
			var err error
			newPostsCompressed[i], err = compressPost(post)
			if err != nil {
				return err
			}
		}

		finalPosts := append(oldDynamoDayFilled.Posts, newPostsCompressed...)
		dayTimeline := &DynamoDayUserTimelinePage{
			PK:         buildPK(userID),
			SK:         buildSK(dayKey),
			Posts:      finalPosts,
			LastUpdate: time.Now(),
			UserID:     userID,
		}

		item, err := attributevalue.MarshalMap(dayTimeline)
		if err != nil {
			fmt.Println(err)
			return err
		}
		transactItem := types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(d.tableName),
				Item:      item,
			},
		}
		transactItems = append(transactItems, transactItem)
	}

	//TODO Validte len because dynamo has 100 itemes limits
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	}

	_, err := d.client.TransactWriteItems(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (d *DynamoDayTimelineFilledRepository) UpdatePosts(ctx context.Context, userID string, post *posts.Post) error {
	dayKey := buildDateKeyByPost(*post)
	dayTimeline, err := d.getDayFilled(ctx, userID, dayKey)
	if err != nil {
		return err
	}

	newPostCompressed, err := compressPost(*post)
	if err != nil {
		return err
	}

	for i, compressedPost := range dayTimeline.Posts {
		decompressedPost, err := decompressPost(compressedPost)
		if err != nil {
			return err
		}

		if decompressedPost.ID == post.ID {
			dayTimeline.Posts[i] = newPostCompressed
			break
		}
	}

	item, err := attributevalue.MarshalMap(dayTimeline)
	if err != nil {
		return err
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	}

	_, err = d.client.PutItem(ctx, putInput)
	return err
}

func (d *DynamoDayTimelineFilledRepository) getDayFilled(ctx context.Context, userID string, dayKey string) (*DynamoDayUserTimelinePage, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: buildPK(userID)},
			"sk": &types.AttributeValueMemberS{Value: buildSK(dayKey)},
		},
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return &DynamoDayUserTimelinePage{
			PK:         buildPK(userID),
			SK:         buildSK(dayKey),
			Posts:      nil,
			LastUpdate: time.Now(),
			UserID:     userID,
		}, nil
	}

	var dayTimeline DynamoDayUserTimelinePage
	err = attributevalue.UnmarshalMap(result.Item, &dayTimeline)
	if err != nil {
		return nil, err
	}
	return &dayTimeline, nil
}

func (d *DynamoDayTimelineFilledRepository) RemovePost(ctx context.Context, userID string, post *posts.Post) error {
	dayKey := buildDateKeyByPost(*post)
	dayTimeline, err := d.getDayFilled(ctx, userID, dayKey)
	if err != nil {
		return err
	}

	var newPosts []string
	postFound := false
	for _, compressedPost := range dayTimeline.Posts {
		decompressedPost, err := decompressPost(compressedPost)
		if err != nil {
			return err
		}

		if decompressedPost.ID == post.ID {
			postFound = true
			continue
		}
		newPosts = append(newPosts, compressedPost)
	}

	if !postFound {
		return nil
	}

	updatedDayTimeline := DynamoDayUserTimelinePage{
		PK:         dayTimeline.PK,
		SK:         dayTimeline.SK,
		Posts:      newPosts,
		LastUpdate: time.Now(),
		UserID:     dayTimeline.UserID,
		Date:       dayTimeline.Date,
	}

	item, err := attributevalue.MarshalMap(updatedDayTimeline)
	if err != nil {
		return err
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	}

	_, err = d.client.PutItem(ctx, putInput)
	return err
}

type DynamoDayUserTimelinePage struct {
	PK string `dynamodbav:"pk"`
	SK string `dynamodbav:"sk"`

	//Fill other
	Posts      []string  `dynamodbav:"posts"`
	LastUpdate time.Time `dynamodbav:"last_update"`
	UserID     string    `dynamodbav:"user_id"`
	Date       time.Time `dynamodbav:"date"`
}

func (d DynamoDayUserTimelinePage) toDomain() (*day_timeline_filled.DayUserTimelineFilled, error) {
	decompressedPosts := make([]posts.Post, len(d.Posts))
	for i, compressedPost := range d.Posts {
		decompressedPost, err := decompressPost(compressedPost)
		if err != nil {
			return nil, err
		}
		decompressedPosts[i] = *decompressedPost
	}
	return &day_timeline_filled.DayUserTimelineFilled{
		LastUpdate: d.LastUpdate,
		Posts:      decompressedPosts,
		UserID:     d.UserID,
	}, nil
}

func compressPost(post posts.Post) (string, error) {
	jsonData, err := json.Marshal(post)
	if err != nil {
		return "", err
	}

	var compressedPost bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedPost)

	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		return "", fmt.Errorf("error compressing: %v", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(compressedPost.Bytes()), nil
}

func decompressPost(compressedPost string) (*posts.Post, error) {
	compressedData, err := base64.StdEncoding.DecodeString(compressedPost)
	if err != nil {
		return nil, fmt.Errorf("error to decoding on base64: %v", err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("error creating decompressor: %v", err)
	}
	defer gzipReader.Close()

	decompressedPost, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("error decompressing post: %v", err)
	}

	var post posts.Post
	err = json.Unmarshal(decompressedPost, &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func splitPostByDate(domainPosts []posts.Post) map[string][]posts.Post {
	postsMapByDay := make(map[string][]posts.Post)
	for _, post := range domainPosts {
		key := buildDateKeyByPost(post)
		if _, ok := postsMapByDay[key]; ok {
			postsMapByDay[key] = append(postsMapByDay[key], post)
			continue
		}
		postsMapByDay[key] = []posts.Post{post}
	}
	return postsMapByDay
}

func buildDateKeyByPost(post posts.Post) string {
	return fmt.Sprintf("%v:%v:%v", post.PublishedAt.Year(), int(post.PublishedAt.Month()), post.PublishedAt.Day())
}

func buildPK(userId string) string {
	return fmt.Sprintf(pkPrefix, userId)
}
func buildSK(dateKey string) string {
	return fmt.Sprintf(skPrefix, dateKey)
}

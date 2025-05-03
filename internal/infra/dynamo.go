package infra

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tufee/desk-reservation-go/internal/domain"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

type Db struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func CreateDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}
	return dynamodb.NewFromConfig(cfg)
}

func InitializeDB(tableName string) Db {
	return Db{
		TableName:      tableName,
		DynamoDbClient: CreateDynamoClient(),
	}
}

func (db Db) FindUserByEmail(ctx context.Context, email string) (*dynamodb.QueryOutput, error) {
	input := dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("EmailIndex"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
		Limit: aws.Int32(1),
	}

	resp, err := db.DynamoDbClient.Query(ctx, &input)
	if err != nil {
		return nil, pkg.NewInternalServerError("failed to query user by email: %w", err)
	}

	return resp, nil
}

func (db Db) SaveUser(ctx context.Context, user domain.User) error {
	item, err := attributevalue.MarshalMap(user)

	if err != nil {
		return pkg.NewInternalServerError("failed to marshal user item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item:      item,
	}

	if _, err = db.DynamoDbClient.PutItem(ctx, input); err != nil {
		return pkg.NewInternalServerError("failed to save user to dynamoDB: %w", err)
	}

	return nil
}

package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ChatDataType struct {
	ChatID    string `json:"chat_id" dynamodbav:"chat_id"`
	UserID    string `json:"user_id" dynamodbav:"user_id"`
	Title     string `json:"title" dynamodbav:"title"`
	CreatedAt int    `json:"created_at" dynamodbav:"created_at"`
}

func Create(client *dynamodb.Client, chatData ChatDataType) (*dynamodb.PutItemOutput, error) {
	av, err := attributevalue.MarshalMap(chatData)
	if err != nil {
		fmt.Printf("Got error marshalling data: %s\n", err)
		return nil, err
	}
	// save chat to db
	output, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Chats"), Item: av,
	})
	if err != nil {
		fmt.Printf("Couldn't add item to table.: %v\n", err)
	}
	return output, err
}

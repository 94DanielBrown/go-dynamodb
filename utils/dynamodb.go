package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func DescribeTable(client *dynamodb.Client, tableName string) (*dynamodb.DescribeTableOutput, error) {
	table, err := client.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
	)

	return table, err
}

func ListTables(client *dynamodb.Client) (*dynamodb.ListTablesOutput, error) {
	tables, err := client.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	)

	return tables, err
}

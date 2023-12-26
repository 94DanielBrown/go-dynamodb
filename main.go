package main

import (
	"context"
	"fmt"
	"github.com/94danielbrown/go-dynamodb/infrastructure"
	"github.com/94danielbrown/go-dynamodb/initializers"
	"github.com/94danielbrown/go-dynamodb/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"time"
)

func main() {
	initializers.LoadEnvVariables()

	config, err := infrastructure.NewAwsConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := infrastructure.NewDynamoDBClient(config)

	s3Client := s3.NewFromConfig(config)
	bucketName := "cf-templates-15u92x5udtnvy-eu-west-1"
	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("Error listing objects:", err)
		return
	}

	fmt.Println("Objects in the bucket:")
	for _, obj := range resp.Contents {
		fmt.Println(*obj.Key)
	}
	table, err := utils.DescribeTable(client, "Messages")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(
		"Table ID: %s \nTable Name: %s\n\n",
		*table.Table.TableId,
		*table.Table.TableName,
	)

	output, err := utils.ListTables(client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func newClient(profile string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return c, nil
}

func createDynamoDBTable(c *dynamodb.Client,
	tableName string, input *dynamodb.CreateTableInput) error {
	var tableDesc *types.TableDescription
	table, err := c.CreateTable(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create table with %v with error: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(c)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Failed to wait on create table %v with error: %v\n", tableName, err)
		}
		tableDesc = table.TableDescription
	}

	fmt.Println(tableDesc)
	return err
}

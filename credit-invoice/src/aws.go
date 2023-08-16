package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewAwsConfig() *aws.Config {
	return &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String("http://localstack:4566"),
	}
}

func NewDynamoDbClient() *dynamodb.DynamoDB {
	sess, err := session.NewSession(NewAwsConfig())

	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		fmt.Println("Connected to DynamoDB.")
	}

	return dynamodb.New(sess)
}

func NewSqsClient() *sqs.SQS {
	sess, err := session.NewSession(NewAwsConfig())

	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		fmt.Println("Connected to SQS.")
	}

	return sqs.New(sess)
}

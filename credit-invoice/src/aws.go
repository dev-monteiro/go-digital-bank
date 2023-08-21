package main

import (
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

func NewDynamoClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(NewAwsConfig())

	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}

func NewSqsClient() (*sqs.SQS, error) {
	sess, err := session.NewSession(NewAwsConfig())

	if err != nil {
		return nil, err
	}

	return sqs.New(sess), nil
}

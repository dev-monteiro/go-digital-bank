package configuration

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewDynamoClnt() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(newAwsConfig())

	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}

func NewSqsClnt() (*sqs.SQS, error) {
	sess, err := session.NewSession(newAwsConfig())

	if err != nil {
		return nil, err
	}

	return sqs.New(sess), nil
}

func newAwsConfig() *aws.Config {
	return &aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_LOGIN"), os.Getenv("AWS_PASS"), ""),
		Endpoint:    aws.String(os.Getenv("AWS_ENDPOINT")),
	}
}

package connector

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type DynamoConn interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

func NewDynamoConn() (DynamoConn, error) {
	sess, err := session.NewSession(newAwsConfig())

	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}

type SqsConn interface {
	GetQueueUrl(input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error)
	ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

func NewSqsConn() (SqsConn, error) {
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

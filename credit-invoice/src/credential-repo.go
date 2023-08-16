package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CredentialRepo struct {
	dynamoDB *dynamodb.DynamoDB
}

func NewCredentialRepo(dynamoDB *dynamodb.DynamoDB) CredentialRepo {
	return CredentialRepo{dynamoDB: dynamoDB}
}

func (repo *CredentialRepo) getCreditAccountId(customerId string) (int, error) {
	output, err := repo.dynamoDB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("customer-credentials-table"),
		Key: map[string]*dynamodb.AttributeValue{
			"customerId": {
				S: aws.String(customerId),
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if output.Item == nil {
		fmt.Println("Not found.")
		return 0, errors.New("Not found.")
	}

	credential := CustomerCredential{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &credential)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return credential.CreditAccountId, nil
}

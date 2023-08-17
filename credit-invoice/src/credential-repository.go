package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CredentialRepository struct {
	dynamoCli *dynamodb.DynamoDB
}

func NewCredentialRepository(dynamoCli *dynamodb.DynamoDB) CredentialRepository {
	return CredentialRepository{dynamoCli: dynamoCli}
}

func (repo *CredentialRepository) getCreditAccountId(customerId string) (int, error) {
	output, err := repo.dynamoCli.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("customer-credentials-table"),
		Key: map[string]*dynamodb.AttributeValue{
			"CustomerId": {
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

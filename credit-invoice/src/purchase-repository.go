package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PurchaseRepository struct {
	dynamoCli *dynamodb.DynamoDB
}

func NewPurchaseRepository(dynamoCli *dynamodb.DynamoDB) *PurchaseRepository {
	return &PurchaseRepository{dynamoCli: dynamoCli}
}

func (repo *PurchaseRepository) save(purchase Purchase) error {
	item, err := dynamodbattribute.MarshalMap(purchase)
	if err != nil {
		return err
	}

	// TODO: add some attribute exists restriction
	_, err = repo.dynamoCli.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("purchases-table"),
	})
	if err != nil {
		return err
	}

	return nil
}

// TODO: verify the not found case
func (repo *PurchaseRepository) findAllByCreditAccountId(creditAccountId int) ([]Purchase, error) {
	result, err := repo.dynamoCli.Query(&dynamodb.QueryInput{
		TableName:              aws.String("purchases-table"),
		KeyConditionExpression: aws.String("#creditAccountId = :creditAccountId"),
		ExpressionAttributeNames: map[string]*string{
			"#creditAccountId": aws.String("CreditAccountId"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":creditAccountId": {
				N: aws.String(strconv.Itoa(creditAccountId)),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var purchases []Purchase
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases, nil
}

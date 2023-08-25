package main

import (
	"devv-monteiro/go-digital-bank/commons"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PurchaseRepo struct {
	dynamoCli *dynamodb.DynamoDB
}

func NewPurchaseRepo(dynamoCli *dynamodb.DynamoDB) *PurchaseRepo {
	return &PurchaseRepo{dynamoCli: dynamoCli}
}

func (repo *PurchaseRepo) save(purchase commons.PurchaseEvent) error {
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
func (repo *PurchaseRepo) findAllByCreditAccountId(creditAccountId int) ([]commons.PurchaseEvent, error) {
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

	var purchases []commons.PurchaseEvent
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases, nil
}

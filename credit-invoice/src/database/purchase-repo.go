package database

import (
	"devv-monteiro/go-digital-bank/commons"
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	"net/http"
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

func (repo *PurchaseRepo) Save(purchase commons.PurchaseEvent) error {
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

func (repo *PurchaseRepo) FindAllByCreditAccountId(creditAccountId int) ([]commons.PurchaseEvent, *conf.AppError) {
	dynaInput := &dynamodb.QueryInput{
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
	}

	dynaOutput, err := repo.dynamoCli.Query(dynaInput)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	var purchases []commons.PurchaseEvent
	err = dynamodbattribute.UnmarshalListOfMaps(dynaOutput.Items, &purchases)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return purchases, nil
}

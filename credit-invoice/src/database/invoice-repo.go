package database

import (
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type InvoiceRepo struct {
	dynaClnt *dynamodb.DynamoDB
}

func NewInvoiceRepo(dynaClnt *dynamodb.DynamoDB) *InvoiceRepo {
	return &InvoiceRepo{dynaClnt: dynaClnt}
}

func (repo *InvoiceRepo) GetId(cbId int) (string, *conf.AppError) {
	fmt.Println("GetId")

	dynaInput := &dynamodb.QueryInput{
		TableName:              aws.String("invoices-table"),
		IndexName:              aws.String("core-bank-index"),
		KeyConditionExpression: aws.String("CoreBankId = :cbId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":cbId": {N: aws.String(strconv.Itoa(cbId))},
		},
	}

	dynaOutput, err := repo.dynaClnt.Query(dynaInput)
	if err != nil {
		return "", &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if *dynaOutput.Count > 1 {
		return "", &conf.AppError{
			Message:    "Unknown error: " + "There are more than one invoice with the same id.",
			StatusCode: http.StatusInternalServerError,
		}
	}

	if *dynaOutput.Count == 0 {
		return "", nil
	}

	invoices := []Invoice{}
	err = dynamodbattribute.UnmarshalListOfMaps(dynaOutput.Items, &invoices)
	if err != nil {
		return "", &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return invoices[0].Id, nil
}

func (repo *InvoiceRepo) Save(invo Invoice) *conf.AppError {
	fmt.Println("Save")

	item, err := dynamodbattribute.MarshalMap(invo)
	if err != nil {
		return &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	_, err = repo.dynaClnt.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("invoices-table"),
	})
	if err != nil {
		return &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

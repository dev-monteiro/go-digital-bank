package database

import (
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CustomerRepo struct {
	dynaClnt *dynamodb.DynamoDB
}

func NewCustomerRepo(dynamoClnt *dynamodb.DynamoDB) *CustomerRepo {
	return &CustomerRepo{dynaClnt: dynamoClnt}
}

func (repo *CustomerRepo) GetCoreBankId(id string) (int, *conf.AppError) {
	fmt.Println("GetCoreBankId")

	dynaInput := dynamodb.GetItemInput{
		TableName: aws.String("customers-table"),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	dynaOutput, err := repo.dynaClnt.GetItem(&dynaInput)
	if err != nil {
		return 0, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if dynaOutput.Item == nil {
		return 0, &conf.AppError{Message: "Customer not found.", StatusCode: http.StatusNotFound}
	}

	cust := Customer{}
	err = dynamodbattribute.UnmarshalMap(dynaOutput.Item, &cust)
	if err != nil {
		return 0, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return cust.CoreBankId, nil
}

package database

import (
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CustomerRepo struct {
	dynaClnt  *dynamodb.DynamoDB
	tableName *string
}

func NewCustomerRepo(dynamoClnt *dynamodb.DynamoDB) *CustomerRepo {
	tableName := aws.String(os.Getenv("AWS_CUSTOMERS_TABLE_NAME"))
	return &CustomerRepo{dynaClnt: dynamoClnt, tableName: tableName}
}

func (repo *CustomerRepo) FindById(id string) (*Customer, *conf.AppError) {
	log.Println("[CustomerRepo] FindById")

	dynaInput := dynamodb.GetItemInput{
		TableName: repo.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	dynaOutput, err := repo.dynaClnt.GetItem(&dynaInput)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if dynaOutput.Item == nil {
		return nil, &conf.AppError{Message: "Customer not found.", StatusCode: http.StatusNotFound}
	}

	cust := Customer{}
	err = dynamodbattribute.UnmarshalMap(dynaOutput.Item, &cust)
	if err != nil {
		return nil, conf.NewUnknownError(err)
	}

	return &cust, nil
}

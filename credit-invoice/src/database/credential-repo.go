package database

import (
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CredentialRepo struct {
	dynaClnt *dynamodb.DynamoDB
}

func NewCredentialRepo(dynamoClnt *dynamodb.DynamoDB) *CredentialRepo {
	return &CredentialRepo{dynaClnt: dynamoClnt}
}

func (repo *CredentialRepo) GetCreditAccountId(customerId string) (int, *conf.AppError) {
	dynaInput := dynamodb.GetItemInput{
		TableName: aws.String("customer-credentials-table"),
		Key: map[string]*dynamodb.AttributeValue{
			"CustomerId": {
				S: aws.String(customerId),
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

	cred := CustomerCredential{}
	err = dynamodbattribute.UnmarshalMap(dynaOutput.Item, &cred)
	if err != nil {
		return 0, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return cred.CreditAccountId, nil
}

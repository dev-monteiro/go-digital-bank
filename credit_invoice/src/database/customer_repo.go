package database

import (
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CustomerRepo interface {
	FindById(id string) (*Customer, error)
}

type customerRepo struct {
	dynaConn  conn.DynamoConn
	tableName *string
}

func NewCustomerRepo(dynaConn conn.DynamoConn) CustomerRepo {
	tableName := aws.String(os.Getenv("AWS_CUSTOMERS_TABLE_NAME"))
	return &customerRepo{dynaConn: dynaConn, tableName: tableName}
}

func (repo *customerRepo) FindById(id string) (*Customer, error) {
	log.Println("[CustomerRepo] FindById")

	dynaInput := dynamodb.GetItemInput{
		TableName: repo.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	dynaOutput, err := repo.dynaConn.GetItem(&dynaInput)
	if err != nil {
		return nil, err
	}

	if dynaOutput.Item == nil {
		return nil, nil
	}

	cust := Customer{}
	err = dynamodbattribute.UnmarshalMap(dynaOutput.Item, &cust)
	if err != nil {
		return nil, err
	}

	return &cust, nil
}

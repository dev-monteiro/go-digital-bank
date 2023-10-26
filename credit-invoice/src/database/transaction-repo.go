package database

import (
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type TransactionRepo struct {
	dynaClnt *dynamodb.DynamoDB
}

func NewTransactionRepo(dynaClnt *dynamodb.DynamoDB) *TransactionRepo {
	return &TransactionRepo{dynaClnt: dynaClnt}
}

func (repo *TransactionRepo) Save(trsac Transaction) error {
	log.Println("[TransactionRepo] Save")

	item, err := dynamodbattribute.MarshalMap(trsac)
	if err != nil {
		return err
	}

	// TODO: add some attribute exists restriction
	_, err = repo.dynaClnt.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("transactions-table"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *TransactionRepo) FindAllByCustomerCoreBankId(custCoreBankId int) ([]Transaction, *conf.AppError) {
	log.Println("[TransactionRepo] FindAllByCustomerCoreBankId")

	dynaInput := &dynamodb.QueryInput{
		TableName:              aws.String("transactions-table"),
		IndexName:              aws.String("customerCoreBankId-index"),
		KeyConditionExpression: aws.String("CustomerCoreBankId = :custCoreBankId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":custCoreBankId": {N: aws.String(strconv.Itoa(custCoreBankId))},
		},
	}

	dynaOutput, err := repo.dynaClnt.Query(dynaInput)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	var transactions []Transaction
	err = dynamodbattribute.UnmarshalListOfMaps(dynaOutput.Items, &transactions)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return transactions, nil
}

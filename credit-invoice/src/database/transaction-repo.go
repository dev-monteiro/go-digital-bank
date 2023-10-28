package database

import (
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type TransactionRepo struct {
	dynaClnt  *dynamodb.DynamoDB
	tableName *string
}

func NewTransactionRepo(dynaClnt *dynamodb.DynamoDB) *TransactionRepo {
	tableName := aws.String(os.Getenv("AWS_TRANSACTIONS_TABLE_NAME"))
	return &TransactionRepo{dynaClnt: dynaClnt, tableName: tableName}
}

func (repo *TransactionRepo) Save(transc Transaction) error {
	log.Println("[TransactionRepo] Save")

	item, err := dynamodbattribute.MarshalMap(transc)
	if err != nil {
		return err
	}

	// TODO: add some attribute exists restriction
	_, err = repo.dynaClnt.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: repo.tableName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *TransactionRepo) FindAllByCustomerCoreBankId(custCoreBankId int) ([]Transaction, *conf.AppError) {
	log.Println("[TransactionRepo] FindAllByCustomerCoreBankId")

	dynaInput := &dynamodb.QueryInput{
		TableName:              repo.tableName,
		IndexName:              aws.String("customerCoreBankId-index"),
		KeyConditionExpression: aws.String("CustomerCoreBankId = :custCoreBankId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":custCoreBankId": {N: aws.String(strconv.Itoa(custCoreBankId))},
		},
	}

	dynaOutput, err := repo.dynaClnt.Query(dynaInput)
	if err != nil {
		return nil, conf.NewUnknownError(err)
	}

	var transactions []Transaction
	err = dynamodbattribute.UnmarshalListOfMaps(dynaOutput.Items, &transactions)
	if err != nil {
		return nil, conf.NewUnknownError(err)
	}

	return transactions, nil
}

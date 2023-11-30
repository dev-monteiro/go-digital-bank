package database

import (
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type TransactionRepo interface {
	Save(transc *Transaction) error
	FindAllByCustomerCoreBankId(custCoreBankId int) ([]*Transaction, error)
	Delete(transc *Transaction) error
}

type transactionRepo struct {
	dynaConn  conn.DynamoConn
	tableName *string
}

func NewTransactionRepo(dynaConn conn.DynamoConn) TransactionRepo {
	tableName := aws.String(os.Getenv("AWS_TRANSACTIONS_TABLE_NAME"))
	return &transactionRepo{dynaConn: dynaConn, tableName: tableName}
}

func (repo *transactionRepo) Save(transc *Transaction) error {
	log.Println("[TransactionRepo] Save")

	item, err := dynamodbattribute.MarshalMap(transc)
	if err != nil {
		return err
	}

	_, err = repo.dynaConn.PutItem(&dynamodb.PutItemInput{
		Item:                item,
		TableName:           repo.tableName,
		ConditionExpression: aws.String("attribute_not_exists(PurchaseId)"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *transactionRepo) FindAllByCustomerCoreBankId(custCoreBankId int) ([]*Transaction, error) {
	log.Println("[TransactionRepo] FindAllByCustomerCoreBankId")

	dynaInput := &dynamodb.QueryInput{
		TableName:              repo.tableName,
		IndexName:              aws.String("customerCoreBankId-index"),
		KeyConditionExpression: aws.String("CustomerCoreBankId = :custCoreBankId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":custCoreBankId": {N: aws.String(strconv.Itoa(custCoreBankId))},
		},
	}

	dynaOutput, err := repo.dynaConn.Query(dynaInput)
	if err != nil {
		return nil, err
	}

	var transactions []*Transaction
	err = dynamodbattribute.UnmarshalListOfMaps(dynaOutput.Items, &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *transactionRepo) Delete(transc *Transaction) error {
	log.Println("[TransactionRepo] Delete")

	input := &dynamodb.DeleteItemInput{
		TableName: repo.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"PurchaseId": {
				N: aws.String(strconv.Itoa(transc.PurchaseId)),
			},
		},
	}

	_, err := repo.dynaConn.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}

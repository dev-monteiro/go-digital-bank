#!/bin/bash

awslocal sqs create-queue --queue-name purchases-queue
awslocal dynamodb create-table --table-name customer-credentials-table --attribute-definitions AttributeName=CustomerId,AttributeType=S --key-schema AttributeName=CustomerId,KeyType=HASH --billing-mode PAY_PER_REQUEST
awslocal dynamodb create-table --table-name invoice-entries-table --attribute-definitions AttributeName=CreditAccountId,AttributeType=N --key-schema AttributeName=CreditAccountId,KeyType=HASH --billing-mode PAY_PER_REQUEST
awslocal dynamodb create-table --table-name purchases-table --attribute-definitions AttributeName=CreditAccountId,AttributeType=N AttributeName=PurchaseId,AttributeType=N --key-schema AttributeName=CreditAccountId,KeyType=HASH AttributeName=PurchaseId,KeyType=RANGE --billing-mode PAY_PER_REQUEST

awslocal dynamodb put-item --table-name customer-credentials-table --item '{"CustomerId": {"S": "abc-123-def"}, "CreditAccountId": {"N": "123"}}'
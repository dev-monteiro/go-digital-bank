#!/bin/bash

awslocal dynamodb create-table --table-name customer-credentials-table --attribute-definitions AttributeName=customerId,AttributeType=S --key-schema AttributeName=customerId,KeyType=HASH --billing-mode PAY_PER_REQUEST
awslocal dynamodb put-item --table-name customer-credentials-table --item '{"customerId": {"S": "abc-123-def"}, "creditAccountId": {"N": "123"}}'

awslocal dynamodb create-table --table-name invoice-entries-table --attribute-definitions AttributeName=creditAccountId,AttributeType=N --key-schema AttributeName=creditAccountId,KeyType=HASH --billing-mode PAY_PER_REQUEST

awslocal sqs create-queue --queue-name purchases-queue
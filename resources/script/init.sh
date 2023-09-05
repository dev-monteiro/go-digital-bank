#!/bin/bash

awslocal sqs create-queue --queue-name purchases-queue \
    > /dev/null

awslocal dynamodb create-table --table-name customers-table \
    --attribute-definitions AttributeName=Id,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    > /dev/null
awslocal dynamodb create-table --table-name invoices-table \
    --attribute-definitions AttributeName=Id,AttributeType=S AttributeName=CoreBankId,AttributeType=N \
    --key-schema AttributeName=Id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes "IndexName=core-bank-index,KeySchema=[{AttributeName=CoreBankId,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
    > /dev/null
awslocal dynamodb create-table --table-name purchases-table \
    --attribute-definitions AttributeName=CreditAccountId,AttributeType=N AttributeName=PurchaseId,AttributeType=N \
    --key-schema AttributeName=CreditAccountId,KeyType=HASH AttributeName=PurchaseId,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    > /dev/null

awslocal dynamodb put-item --table-name customers-table \
    --item '{"Id": {"S": "abc-123-def"}, "CoreBankId": {"N": "123"}}' \
    > /dev/null

echo "Setup completed"
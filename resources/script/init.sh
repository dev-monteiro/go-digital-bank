#!/bin/bash

awslocal dynamodb create-table --table-name customers-table \
    --attribute-definitions AttributeName=Id,AttributeType=S AttributeName=CoreBankBatchId,AttributeType=N \
    --key-schema AttributeName=Id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes "IndexName=coreBankBatch-index,KeySchema=[{AttributeName=CoreBankBatchId,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
    > /dev/null

awslocal dynamodb create-table --table-name transactions-table \
    --attribute-definitions AttributeName=PurchaseId,AttributeType=N AttributeName=CustomerCoreBankId,AttributeType=N \
    --key-schema AttributeName=PurchaseId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes "IndexName=customerCoreBankId-index,KeySchema=[{AttributeName=CustomerCoreBankId,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
    > /dev/null

awslocal dynamodb put-item --table-name customers-table \
    --item '{"Id": {"S": "abc-123-def"}, "CoreBankId": {"N": "123"}, "CoreBankBatchId": {"N": "789"}}' \
    > /dev/null

awslocal sqs create-queue --queue-name purchases-queue \
   > /dev/null

echo "setup completed!"
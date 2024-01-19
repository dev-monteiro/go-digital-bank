#!/bin/bash

docker exec -i localstack aws configure set aws_access_key_id test --profile default
docker exec -i localstack aws configure set aws_secret_access_key test --profile default
#docker exec -i localstack rm ~/.aws/config

docker exec -i localstack aws dynamodb scan --endpoint-url http://localhost:4566 --region us-east-1 --table-name local-customers-table --projection-expression "Id" --output json \
        | sed -n 's|.*"Id":{"S":"\([^"]*\)".*|\1|p' \
        | while read -r key; do
            docker exec -i localstack aws dynamodb delete-item --endpoint-url http://localhost:4566 --region us-east-1 --table-name local-customers-table --key "{\"Id\":{\"S\":\"$key\"}}"
        done

docker exec -i localstack aws dynamodb scan --endpoint-url http://localhost:4566 --region us-east-1 --table-name local-transactions-table --projection-expression "PurchaseId" --output json \
        | sed -n 's|.*"PurchaseId":{"S":"\([^"]*\)".*|\1|p' \
        | while read -r key; do
            docker exec -i localstack aws dynamodb delete-item --endpoint-url http://localhost:4566 --region us-east-1 --table-name local-customers-table --key "{\"PurchaseId\":{\"S\":\"$key\"}}"
        done
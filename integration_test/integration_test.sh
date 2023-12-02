#!/bin/batch

cd ..
docker-compose -f docker-compose.test.yml -f docker-compose.yml up --build -d

cd integration_test
go test

cd ..
docker-compose down
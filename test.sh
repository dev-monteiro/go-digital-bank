#!/bin/batch

docker-compose -f docker-compose.test.yml -f docker-compose.yml up --build --remove-orphans -d

cd integration_test
go test

cd ..
docker-compose down
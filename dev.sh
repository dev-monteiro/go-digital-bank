#!/bin/batch

docker-compose -f docker-compose.dev.yml -f docker-compose.yml up --build --remove-orphans
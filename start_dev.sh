#! /usr/bin/env bash

export PORT=8080
export DB_USERNAME=web
export DB_PASSWORD=web
export DB_HOST=localhost
export DB_NAME=web

go run cmd/main.go

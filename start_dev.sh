#! /usr/bin/env bash

export PORT=8080
export DB_USERNAME=web
export DB_PASSWORD=web
export DB_HOST=localhost

go run cmd/main.go

#! /usr/bin/env bash

export PORT=8080
export DB_USERNAME=web
export DB_PASSWORD=web
export DB_HOST=localhost
export DB_NAME=web
export GOOGLE_CLIENT_ID=764012067069-qqcmc5c5kupij0vho6ue312el2ujn52l.apps.googleusercontent.com
export GOOGLE_CLIENT_SECRET=_268NyCGn_F0q4WUM9P35jIh
export GOOGLE_WHITE_LIST_EMAILS=rcliao01@gmail.com,jason@may54th.com
export GOOGLE_REDIRECT_URI=http://localhost:8080/google/callback

go run cmd/main.go

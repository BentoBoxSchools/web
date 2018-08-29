#! /usr/bin/env bash

docker run --name bentobox-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:5
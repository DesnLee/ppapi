#! /usr/bin/env bash

sqlc generate && swag fmt && swag init --parseDependency && go build -o ppapi && ./ppapi server

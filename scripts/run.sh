#! /usr/bin/env bash

sqlc generate && swag fmt && swag init && go build -o ppapi && ./ppapi server

#! /usr/bin/env bash

sqlc generate && swag init && go build -o ppapi && ./ppapi server

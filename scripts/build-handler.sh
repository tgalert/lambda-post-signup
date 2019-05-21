#!/bin/bash

GOOS=linux GOARCH=amd64 go build src/handler.go src/aws.go src/k8s.go src/http.go

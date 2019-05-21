#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o assets/aws-iam-authenticator github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator

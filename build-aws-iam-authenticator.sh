#!/bin/bash

GOOS=linux GOARCH=amd64 go build github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator

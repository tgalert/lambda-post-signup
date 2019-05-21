#!/bin/bash

physical_resource_id=$(aws cloudformation describe-stack-resources \
    --stack-name tgalert-lambda-post-signup \
    --query "StackResources[?ResourceType == 'AWS::IAM::Role'].PhysicalResourceId" \
    --output text)

arn=$(aws iam get-role \
    --role-name "$physical_resource_id" \
    --query "Role.Arn" \
    --output text)

echo "$arn"

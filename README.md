# TG-Search Post Signup Lambda Function

This is an AWS Lambda application consisting of a single Lambda function that is invoked whenever a new user signs up to TG-Search.

## Actions

This Lambda function creates the following resources for the new TG-Search user:

- A new RabbitMQ user and vhost in the RabbitMQ server
- A new [tgsearch-core](https://hub.docker.com/r/weibeld/tg-monitor/) Kubernetes deployment

## Triggering

To trigger the Lambda function, we make use of Amazon Cognito [User Pool Lambda Triggers](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html).

In particular, we use a [post confirmation Lambda trigger](https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-post-confirmation.html). This trigger is invoked whenever the user correctly enters a confirmation code.

There are two scenarios in which the user needs to enter a confirmation code. First, to confirm the initial sign-up, and second, to confirm a password reset request ("forgot password"). Each event has its own [trigger source](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html#cognito-user-identity-pools-working-with-aws-lambda-trigger-sources) value that is included in the event data that is sent to the Lambda function:

- Sign-up confirmation: `PostConfirmation_ConfirmSignUp`
- Forgot password confirmation: `PostConfirmation_ConfirmForgotPassword`

Our Lambda function executes the above actions only for `PostConfirmation_ConfirmSignUp` events.

## Deployment

### Deploy Lambda Application

The Lambda application can be deployed with the [AWS SAM CLI](https://github.com/awslabs/aws-sam-cli):

~~~bash
sam package \
  --template-file template.yml \
  --output-template-file package.yml \
  --s3-bucket quantumsense-sam
  
sam deploy \
  --template-file package.yml \
  --stack-name tgsearch-post-signup \
  --capabilities CAPABILITY_IAM
~~~

### Set Up Cognito User Pool Lambda Trigger

If it's a new deployment, you have to set the Lambda trigger for the user pool in the [AWS Console](https://console.aws.amazon.com/cognito/users).

## Local Testing

The Lambda function can be invoked locally with the SAM CLI and some sample event data:

~~~bash
sam local invoke --event event.json PostSignupFunction
~~~

The [event.json](event.json) file contains a real-world sample Cognito User Pool Lambda Trigger event that can be used for local testing:

~~~json
{
    "version": "1",
    "region": "us-east-1",
    "userPoolId": "us-east-1_caEtzyjHJ",
    "userName": "danielmweibel@gmail.com",
    "callerContext": {
        "awsSdkVersion": "aws-sdk-unknown-unknown",
        "clientId": "dcvoqr2l6ibq6hn6rnlt66man"
    },
    "triggerSource": "PostConfirmation_ConfirmSignUp",
    "request": {
        "userAttributes": {
            "sub": "71dd06a0-e7c5-40c5-a538-3f9a54c947d6",
            "cognito:user_status": "CONFIRMED",
            "email_verified": "true",
            "email": "danielmweibel@gmail.com"
        }
    },
    "response": {}
}
~~~

# Lambda Execution Role

AWS SAM creates a default role for each Lambda function. This role has only the following policy:

- [arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole](https://console.aws.amazon.com/iam/home?#/policies/arn:-aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole$jsonEditor)

This policy allows only to create and write CloudWatch logs.

If you use this default role and your Lambda function tries to access another AWS service, you get a runtime error like this:

~~~
AccessDeniedException: User: arn:aws:sts::202449302273:assumed-role/tgalert-lambda-post-signup-PostSignupFunctionRole-1LHTYVZDD0F93/tgalert-lambda-post-signup-PostSignupFunction-T2877QBTM3RJ is not authorized to perform: secretsmanager:GetSecretValue on resource: arn:aws:secretsmanager:us-east-1:202449302273:secret:tgalert/rabbitmq-admin-dh2yZp
	status code: 400, request id: bdc8abd8-1cbb-40e7-8a80-63a12c7902b1
~~~

To allow the Lambda function to access this AWS service, we need to either extend the default role with new [policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html) or define a new role from the ground up and associate it with the Lambda function. SAM provides support for both cases.

## Extend default role

Use the **Policies** property of the [AWS::Serverless::Function](https://github.com/awslabs/serverless-application-model/blob/develop/versions/2016-10-31.md#awsserverlessfunction) to define policies that will be added to the default role:

~~~yaml
Resources:
  PostSignupFunction:
    Properties:
      Policies:
        - Version: "2012-10-17"
          Statement: 
            - Effect: "Allow"
              Action: "secretsmanager:GetSecretValue"
              Resource: "arn:aws:secretsmanager:*:*:secret:*"
        - Version: "2012-10-17"
          Statement: 
            - Effect: "Allow"
              Action: "ssm:GetParameter"
              Resource: "arn:aws:ssm:us-east-1:202449302273:parameter//tgalert/rabbitmq-host"
~~~

There are three ways to define policies in the Policies property:

- Managed policy names
- Inline policy definitions
- [SAM policy templates](https://github.com/awslabs/serverless-application-model/blob/develop/docs/policy_templates.rst)

SAM policy templates are a number of useful pre-defined policies for use in SAM templates. All the policy templates are listed [here](https://github.com/awslabs/serverless-application-model/blob/develop/samtranslator/policy_templates_data/policy_templates.json).

## Define new role from scratch

Use the **Role** property of the [AWS::Serverless::Function](https://github.com/awslabs/serverless-application-model/blob/develop/versions/2016-10-31.md#awsserverlessfunction) to define the ARN of the role that you want to assign to the Lambda function.

You can define the role itself in the same SAM template as a CloudFormation [AWS::IAM::Role](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-role.html) resource:

~~~yaml
Resources:

  PostSignupFunction:
    Properties:
      [...]
      Role: !GetAtt PostSignupFunctionRole.Arn

  PostSignupFunctionRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement: 
          - Effect: "Allow"
            Principal: 
              Service: 
                - "lambda.amazonaws.com"
            Action: 
              - "sts:AssumeRole"
      Policies:
        - PolicyName: "SecretsManagerGetSecretValue"
          PolicyDocument:
            Version: "2012-10-17"
            Statement: 
              - Effect: "Allow"
                Action: "secretsmanager:GetSecretValue"
                Resource: "arn:aws:secretsmanager:*:*:secret:*"
        - PolicyName: "SystemsManagerGetParameter"
          PolicyDocument:
            Version: "2012-10-17"
            Statement: 
              - Effect: "Allow"
                Action: "ssm:GetParameter"
                Resource: "*"
        - PolicyName: "CloudWatchLogs"
          PolicyDocument:
            Version: "2012-10-17"
            Statement: 
              - Effect: "Allow"
                Resource: "*"
                Action:
                  - "logs:CreateLogGroup"
                  - "logs:CreateLogStream"
                  - "logs:PutLogEvents"
~~~

You can define the policies for the role in the [Policies](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-role.html#cfn-iam-role-policies) property of the role itself, or you can define them as separate [AWS::IAM::Policy](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-policy.html) resources and attach them to the role in the [Roles](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-policy.html#cfn-iam-policy-roles) property.

You can also create standalon policies that are listed as IAM resources [here](https://console.aws.amazon.com/iam/home?#/policies). These policies are called [managed policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_managed-vs-inline.html). There is a large number of managed policies that managed by AWS. These policies are called **AWS managed policies**. But you can also create your own standalone policies, which then are called **customer managed policies**.

When you define a role, you can also attach these managed policies to them in the [ManagedPolicyArns](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-role.html#cfn-iam-role-managepolicyarns) property.


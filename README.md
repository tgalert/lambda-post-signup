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

It's possible to execute the Lambda function locally with the SAM CLI and some sample event data:

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
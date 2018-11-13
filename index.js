exports.handler = async function(event, context) {

  // This event is also triggered on "forgot password" confirmation, in which
  // case the trigger source is 'PostConfirmation_ConfirmForgotPassword':
  // https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html#cognito-user-identity-pools-working-with-aws-lambda-trigger-sources
  if (event.triggerSource === 'PostConfirmation_ConfirmSignUp') {
    /* Sample event JSON:
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
    */
    console.log(`User ${event.userName} (${event.request.userAttributes.sub})`);
  }
  return event;
};
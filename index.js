exports.handler = async function(event, context) {

  // This event is also triggered on "forgot password" confirmation, in which
  // case the trigger source is 'PostConfirmation_ConfirmForgotPassword':
  // https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html#cognito-user-identity-pools-working-with-aws-lambda-trigger-sources
  if (event.triggerSource === 'PostConfirmation_ConfirmSignUp') {
    // See event.json file for sample event data
    console.log(`User ${event.userName} (${event.request.userAttributes.sub})`);
  }
  return event;
};
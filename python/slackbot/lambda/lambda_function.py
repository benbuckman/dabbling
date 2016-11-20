'''
AWS Lambda function for responding to Slack "Slash command"

http://docs.aws.amazon.com/lambda/latest/dg/python-programming-model-handler-types.html
https://api.slack.com/slash-commands

Based on existing Lambda Slack template.

For generating the encrypted verification key:
    $ aws kms encrypt --profile <name> --key-id alias/<KMS key name> --region us-west-2 --plaintext "<verification>"
'''

import boto3
import json
import logging
import os

from base64 import b64decode
from urlparse import parse_qs

ENCRYPTED_EXPECTED_TOKEN = os.environ['kmsEncryptedToken']

kms = boto3.client('kms')
expected_token = kms.decrypt(CiphertextBlob=b64decode(ENCRYPTED_EXPECTED_TOKEN))['Plaintext']

logger = logging.getLogger()
logger.setLevel(logging.DEBUG)


def build_response(err, msg=None):
    response = {}

    if err:
        response['statusCode'] = '400'
        response['body'] = err.message
    else:
        response['statusCode'] = '200'
        response['body'] = json.dumps(
            {
                "response_type": "in_channel",
                "text": msg
            }
        )

    response['headers'] = {
        'Content-Type': 'application/json'
    }

    logger.info("Responding: %s " % response)
    return response


def respond_to_command(text, user=None, channel=None):
    if text == "help":
        return list_commands()

    elif text == "learn":
        return random_wikipedia_link()

    else:
        return build_response(
            None,
            "Hello %s! I don't understand '%s'." % (user, text)
        )


def list_commands():
    return build_response(None, "Valid commands are: help, learn")


def random_wikipedia_link():
    return build_response(None, "Coming soon!")


def lambda_handler(event, context):
    logger.info("Received: %s" % event)

    params = parse_qs(event['body'])
    token = params['token'][0]
    if token != expected_token:
        logger.error("Request token (%s) does not match expected: %s" % token)
        return build_response(Exception('Invalid request token'))

    try:
        return respond_to_command(
            params['text'][0],
            params['user_name'][0],
            params['channel_name'][0]
        )

    except:
        return build_response(
            None,
            "Unexpected error: %s" % sys.exc_info()[0]
        )

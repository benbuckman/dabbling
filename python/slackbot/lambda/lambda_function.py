'''
AWS Lambda function for responding to Slack "Slash command"
(https://api.slack.com/slash-commands)
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


def respond(err, msg=None):
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


def lambda_handler(event, context):
    logger.info("Received: %s" % event)

    params = parse_qs(event['body'])
    token = params['token'][0]
    if token != expected_token:
        logger.error("Request token (%s) does not match expected: %s" % token)
        return respond(Exception('Invalid request token'))

    user = params['user_name'][0]
    command = params['command'][0]
    channel = params['channel_name'][0]
    command_text = params['text'][0]

    return respond(None, "Hello! %s invoked %s in %s with the following text: %s" % (user, command, channel, command_text))

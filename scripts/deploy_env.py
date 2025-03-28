# Use this code snippet in your app.
# If you need more information about configurations
# or implementing the sample code, visit the AWS docs:
# https://aws.amazon.com/developer/language/python/

import boto3
from botocore.exceptions import ClientError


def get_secret():
    secret_name = "tracker-backend-keys"
    region_name = "us-west-1"

    # Create a Secrets Manager client
    session = boto3.session.Session()
    client = session.client(
        service_name='secretsmanager',
        region_name=region_name
    )

    try:
        get_secret_value_response = client.get_secret_value(
            SecretId=secret_name
        )
    except ClientError as e:
        # For a list of exceptions thrown, see
        # https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
        raise e

    secret1 = get_secret_value_response['B_PORT']
    secret2 = get_secret_value_response['POSTGRES_HOST']
    secret3 = get_secret_value_response['POSTGRES_USER']
    secret4 = get_secret_value_response['POSTGRES_PORT']
    secret5 = get_secret_value_response['POSTGRES_PASSWORD']
    secret6 = get_secret_value_response['POSTGRES_URL']
    secret7 = get_secret_value_response['DB_NAME']
    secret8 = get_secret_value_response['JWT']
    secret9 = get_secret_value_response['ORG_ADD_KEY']
    secret10 = get_secret_value_response['ORG_ADD_KEY_1']
    secret11 = get_secret_value_response['ORG_ADD_KEY_2']
    secret12 = get_secret_value_response['ORG_ADD_KEY_3']

    print(secret1, secret2, secret3)

    # Your code goes here.

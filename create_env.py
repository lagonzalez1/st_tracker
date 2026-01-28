import os
import json
import boto3
from pathlib import Path

class Secrets:
    def __init__(self, data):
        self.B_PORT = data.get("B_PORT")
        self.POSTGRES_HOST = data.get("POSTGRES_HOST")
        self.POSTGRES_USER = data.get("POSTGRES_USER")
        self.POSTGRES_PORT = data.get("POSTGRES_PORT")
        self.POSTGRES_PASSWORD = data.get("POSTGRES_PASSWORD")
        self.POSTGRES_URL = data.get("POSTGRES_URL")
        self.DB_NAME = data.get("DB_NAME")
        self.JWT = data.get("JWT")
        self.ORG_ADD_KEY = data.get("ORG_ADD_KEY")
        self.ORG_ADD_KEY_1 = data.get("ORG_ADD_KEY_1")
        self.ORG_ADD_KEY_2 = data.get("ORG_ADD_KEY_2")
        self.ORG_ADD_KEY_3 = data.get("ORG_ADD_KEY_3")
        self.AMAZON_MQ = data.get("AMAZON_MQ")
        self.RABBIT_USERNAME = data.get("RABBIT_USERNAME")
        self.RABBIT_PASSWORD = data.get("RABBIT_PASSWORD")
        self.VALKEY_PORT = data.get("VALKEY_PORT")
        self.VALKEY_HOST = data.get("VALKEY_HOST")
        self.PROD_URL = data.get("PROD_URL")
        self.DEV_URL = data.get("DEV_URL")
        self.STRIPE_SECRET = data.get("STRIPE_SECRET")
        self.STRIPE_PUBLIC = data.get("STRIPE_PUBLIC")
        self.STRIPE_WEBHOOK = data.get("STRIPE_WEBHOOK")
        self.PROD_QUEUE_NAME_DATA_REPORTS = data.get("PROD_QUEUE_NAME_DATA_REPORTS")
        self.PROD_QUEUE_NAME_DATA_ASSESSMENTS = data.get("PROD_QUEUE_NAME_DATA_ASSESSMENTS")

    def to_env_dict(self):
        return {
            "DB_HOST": self.B_PORT,
            "POSTGRES_HOST": self.POSTGRES_HOST,
            "POSTGRES_USER": self.POSTGRES_USER,
            "POSTGRES_PORT": self.POSTGRES_PORT,
            "POSTGRES_PASSWORD": self.POSTGRES_PASSWORD,
            "POSTGRES_URL": self.POSTGRES_URL,
            "DB_NAME": self.DB_NAME,
            "JWT": self.JWT,
            "ORG_ADD_KEY": self.ORG_ADD_KEY,
            "ORG_ADD_KEY_1": self.ORG_ADD_KEY_1,
            "ORG_ADD_KEY_2": self.ORG_ADD_KEY_2,
            "ORG_ADD_KEY_3": self.ORG_ADD_KEY_3,
            "AMAZON_MQ": self.AMAZON_MQ,
            "RABBIT_PASSWORD": self.RABBIT_PASSWORD,
            "RABBIT_USERNAME": self.RABBIT_USERNAME,
            "VALKEY_PORT": self.VALKEY_PORT,
            "VALKEY_HOST": self.VALKEY_HOST,
            "PROD_URL": self.PROD_URL,
            "DEV_URL": self.DEV_URL,
            "STRIPE_SECRET": self.STRIPE_SECRET,
            "STRIPE_PUBLIC": self.STRIPE_PUBLIC,
            "STRIPE_WEBHOOK": self.STRIPE_WEBHOOK,
            "PROD_QUEUE_NAME_DATA_REPORTS": self.PROD_QUEUE_NAME_DATA_REPORTS,
            "PROD_QUEUE_NAME_DATA_ASSESSMENTS": self.PROD_QUEUE_NAME_DATA_ASSESSMENTS
        }

## Returns dictionary of env keys.
## Throws exception
def get_aws_secrets():
    secret_name = "tracker-backend-keys"
    region_name = "us-west-1"
    session = boto3.session.Session()
    client = session.client(
        service_name="secretsmanager",
        region_name=region_name
    )
    try:
        response = client.get_secret_value(SecretId=secret_name)
    except Exception as e:
        raise RuntimeError(f"Error retrieving secret: {e}")
    secret_string = response.get("SecretString")
    if not secret_string:
        raise ValueError("Secret string is empty")

    return json.loads(secret_string)


## Check filepath given filepath.
## Returns Boolean
def file_exists(filepath):
    return Path(filepath).is_file()

## Write to env file given file name and dict
## Returns None
def write_env_file(filename, secrets_dict):
    try:
        with open(filename, "w") as f:
            for key, value in secrets_dict.items():
                if value is not None:
                    f.write(f"{key}={value}\n")
        print("Successfully created", filename)
    except Exception as e:
        raise RuntimeError(f"Unable to write to file: {e}")


## Check .env.production file exist, skip if true otherwise run script.
## Return None
def aws_config():
    secrets_data = get_aws_secrets()
    # Please fix this, no need to parse each value out.
    secrets = Secrets(secrets_data)
    env_path = "/home/ubuntu/app/.env.production"
    write_env_file(env_path, secrets.to_env_dict())

# Run this to trigger the config fetch and file write
if __name__ == "__main__":
    aws_config()

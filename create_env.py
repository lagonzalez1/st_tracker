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
        }

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

def file_exists(filepath):
    return Path(filepath).is_file()

def write_env_file(filename, secrets_dict):
    try:
        with open(filename, "w") as f:
            for key, value in secrets_dict.items():
                if value is not None:
                    f.write(f"{key}={value}\n")
        print("Successfully created", filename)
    except Exception as e:
        raise RuntimeError(f"Unable to write to file: {e}")

def aws_config():
    secrets_data = get_aws_secrets()
    # Please fix this, no need to parse each value out.
    secrets = Secrets(secrets_data)

    env = os.getenv("APP_ENV", "development")
    env_path = "/home/ubuntu/app/.env.production"

    if env == "production" and file_exists(env_path):
        print("File exists...")
        return

    write_env_file(".env.production", secrets.to_env_dict())

# Run this to trigger the config fetch and file write
if __name__ == "__main__":
    aws_config()

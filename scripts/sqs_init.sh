#!/bin/bash
set -e

awslocal sqs create-queue --queue-name dev-data-assessments --region us-west-1

awslocal sqs create-queue --queue-name dev-data-reports --region us-west-1

awslocal sqs create-queue --queue-name dev-generate-generate --region us-west-1

awslocal sqs create-queue --queue-name dev-survey-analysis --region us-west-1

echo "SQS queues created successfully!"

# List all queues to verify
awslocal sqs list-queues
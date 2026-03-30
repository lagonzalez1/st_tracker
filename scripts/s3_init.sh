#!/bin/bash
set -e

awslocal s3 mb s3://tracker-client-storage-dev
echo "S3 Created bucket!"

awslocal s3 mb s3://content-generate-dev
echo "S3 Created bucket!"

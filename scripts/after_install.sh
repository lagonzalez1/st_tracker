#!/bin/bash
set -e

echo "Extracting application..."
cd /home/ubuntu/app
tar -xzvf app.tar.gz
chmod +x app
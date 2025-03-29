#!/bin/bash
set -e

echo "Extracting application..."
# Enter into app folder 
cd /home/ubuntu/app
ls -l 
# Ensure its executable
sudo chmod +x app.tar.gz
# Extract all file
sudo tar -xzf app.tar.gz
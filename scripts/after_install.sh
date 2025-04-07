#!/bin/bash
set -e

echo "Extracting application..."
# Enter into app folder 
cd /home/ubuntu/app
sudo python3 create_env.py



ls -l 
# Ensure its executable


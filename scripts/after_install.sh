#!/bin/bash
set -e

echo "Extracting application..."
# Enter into app folder 
cd /home/ubuntu/app

if [ -f ".env.production"]; then
    echo "Env production exist continue.."
else
    echo "Env production does not exist running to create"
    sudo python3 create_env.py
fi



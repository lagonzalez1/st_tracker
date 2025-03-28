#!/bin/bash
set -e

echo "Starting application..."
cd /home/ubuntu/app

# Ensure old instance is stopped before starting a new one
if [ -f app.pid ]; then
  kill -9 $(cat app.pid) || true
  rm -f app.pid
fi

# Run the compiled Go binary
nohup ./app > app.log 2>&1 &

# Store the process ID for later use
echo $! > app.pid

echo "Application started successfully!"

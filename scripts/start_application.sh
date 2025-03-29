#!/bin/bash
set -e

echo "Starting application..."
pwd
ls -l
# Ensure old instance is stopped before starting a new one
if [ -f app.pid ]; then
  kill -9 $(cat app1.pid) || true
  rm -f app1.pid
fi

# Run the compiled Go binary this might be called ./main
nohup ./app1 > app1.log 2>&1 &

# Store the process ID for later use
echo $! > app1.pid

echo "Application started successfully!"

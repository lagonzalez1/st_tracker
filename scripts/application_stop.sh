#!/bin/bash
set -e

if [ -f /home/ubuntu/app/app.pid ]; then
  echo "Stopping application..."
  kill -9 $(cat /home/ubuntu/app/app.pid) || true
  rm -f /home/ubuntu/app/app.pid
fi
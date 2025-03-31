#!/bin/bash
set -euo pipefail

# Configuration
APP_NAME="app1"
APP_DIR="/home/ubuntu/app"
LOG_FILE="${APP_DIR}/${APP_NAME}.log"
PID_FILE="${APP_DIR}/${APP_NAME}.pid"
MAX_RETRIES=3
RETRY_DELAY=5

echo "Starting application deployment..."
echo "Current directory: $(pwd)"

cd "${APP_DIR}" || {
  echo "Failed to change directory to ${APP_DIR}"
  exit 1
}

# Set production environment
export APP_ENV=production
export GIN_MODE=release  # If using Gin web framework

# Function to stop existing application
stop_application() {
  if [ -f "${PID_FILE}" ]; then
    echo "Found existing PID file. Stopping process..."
    local pid=$(cat "${PID_FILE}")
    
    if ps -p "${pid}" > /dev/null; then
      echo "Stopping process ${pid}..."
      kill -TERM "${pid}" || kill -9 "${pid}" 2>/dev/null
      sleep 2  # Give it time to shut down gracefully
    fi
    
    rm -f "${PID_FILE}"
    echo "Application stopped."
  fi
}

# Function to start application
start_application() {
  echo "Starting ${APP_NAME}..."
  
  # Rotate logs if they exist
  [ -f "${LOG_FILE}" ] && mv "${LOG_FILE}" "${LOG_FILE}.old"
  
  nohup "./${APP_NAME}" >> "${LOG_FILE}" 2>&1 &
  local pid=$!
  
  # Verify process started
  if ! ps -p "${pid}" > /dev/null; then
    echo "Failed to start application!"
    return 1
  fi
  
  echo "${pid}" > "${PID_FILE}"
  echo "Application started with PID ${pid}"
  return 0
}

# Main execution
stop_application

retry_count=0
while [ ${retry_count} -lt ${MAX_RETRIES} ]; do
  if start_application; then
    echo "Application started successfully!"
    exit 0
  fi
  
  retry_count=$((retry_count + 1))
  echo "Attempt ${retry_count} failed. Retrying in ${RETRY_DELAY} seconds..."
  sleep ${RETRY_DELAY}
done

echo "Failed to start application after ${MAX_RETRIES} attempts!"
exit 1
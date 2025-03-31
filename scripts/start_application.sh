#!/bin/bash
set -euo pipefail

# Configuration
APP_NAME="app1"
APP_DIR="/home/ubuntu/app"
cd "${APP_DIR}"

nohup "sudo APP_ENV=production ./${APP_NAME}" >> "${LOG_FILE}" 2>&1 &


LOG_FILE="${APP_DIR}/${APP_NAME}.log"
PID_FILE="${APP_DIR}/${APP_NAME}.pid"

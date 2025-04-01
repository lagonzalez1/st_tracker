#!/bin/bash
set -euo pipefail

# Configuration
APP_NAME="app1"
APP_DIR="/home/ubuntu/app"
SERVICE_FILE="application.service"
SERVICE_NAME="myapp"

echo "Changing to application directory: ${APP_DIR}"
cd "${APP_DIR}" || { echo "Failed to change directory to ${APP_DIR}"; exit 1; }

if [ -f "${SERVICE_FILE}" ]; then
    echo "Starting application deployment"
    
    # Copy service file to systemd directory
    echo "Copying service file to /etc/systemd/system/"
    sudo cp "${SERVICE_FILE}" "/etc/systemd/system/${SERVICE_NAME}.service"
    
    # Set environment variables (correct syntax)
    echo "Setting up environment"
    sudo mkdir -p /etc/systemd/system/"${SERVICE_NAME}".service.d/
    echo -e "[Service]\nEnvironment=APP_ENV=production" | sudo tee /etc/systemd/system/"${SERVICE_NAME}".service.d/override.conf
    
    # Reload and enable service
    echo "Reloading systemd"
    sudo systemctl daemon-reload
    
    echo "Enabling and starting service"
    sudo systemctl enable "${SERVICE_NAME}"
    sudo systemctl restart "${SERVICE_NAME}"  # Using restart instead of start for idempotency
    
    echo "Service status:"
    sudo systemctl status "${SERVICE_NAME}"
else
    echo "Error: Service file ${SERVICE_FILE} not found in ${APP_DIR}"
    exit 1
fi
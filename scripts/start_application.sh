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
    sudo mv application.service /etc/systemd/system
    
    # Reload and enable service
    echo "Reloading systemd"
    sudo systemctl daemon-reload
    echo "Enabling and starting service"
    sudo systemctl enable application
    sudo systemctl start application
    sudo systemctl restart application
    echo "Service status:"
    sudo systemctl status application
else
    echo "Error: Service file ${SERVICE_FILE} not found in ${APP_DIR}"
    exit 1
fi
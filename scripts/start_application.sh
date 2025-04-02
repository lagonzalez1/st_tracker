#!/bin/bash
set -euo pipefail

# Configuration
APP_NAME="app1"
APP_DIR="/home/ubuntu/app"
SERVICE_FILE="application.service"
SERVICE_NAME="myapp"
PORT=3333

echo "Changing to application directory: ${APP_DIR}"
cd "${APP_DIR}" || { echo "Failed to change directory to ${APP_DIR}"; exit 1; }


# Check if port is active
if sudo lsof -i :$PORT >/dev/null; then
    echo "Port $PORT is in use."
    sudo kill -9 $(sudo lsof -t -i:$PORT)
else
    echo "Port $PORT is not in use."
fi

if [ -f "${SERVICE_FILE}" ]; then
    echo "Starting application deployment"
    
    # Copy service file to systemd directory
    echo "Copying service file sudo mv application.service /etc/systemd/system/"
    sudo mv application.service /etc/systemd/system/

    # Reload and enable service 
    echo "Reloading systemd"
    sudo systemctl daemon-reload
    echo "Enabling and starting service"
    sudo systemctl enable application
    sudo systemctl start application

    echo "Service status:"
    sudo systemctl status application
else
    echo "Error: Service file ${SERVICE_FILE} not found in ${APP_DIR}"
    exit 1
fi
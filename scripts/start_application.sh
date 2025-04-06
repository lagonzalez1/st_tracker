#!/bin/bash
set -euo pipefail
set -x 
LOGFILE="/var/log/app-deploy.log"
exec > >(tee -a "$LOGFILE") 2>&1
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
    if [ $? -ne 0 ]; then
        echo "Error: Failed to reload systemd"
        exit 1
    fi

    sudo systemctl enable application
    if [ $? -ne 0 ]; then
        echo "Error: Failed to enable service"
        exit 1
    fi

    sudo systemctl start application
    sleep 2  # Allow service startup time
else
    echo "Error: Service file ${SERVICE_FILE} not found in ${APP_DIR}"
    exit 1
fi


if ! systemctl is-active --quiet application; then
    echo "Error: Service failed to start"
    journalctl -u application -n 20 --no-pager
    exit 1
fi

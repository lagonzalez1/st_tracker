#!/bin/bash
set -euo pipefail

# Configuration
APP_NAME="app1"
APP_DIR="/home/ubuntu/app"
cd "${APP_DIR}"


if [ -f "application.service"]; then
    sudo systemctl daemon-reload
    sudo systemctl start myapp
    sudo systemctl enable myapp 
    sudo systemctl status myapp
fi
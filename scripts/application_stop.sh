#!/bin/bash
set -e
pwd

if [ -f /home/systemd/system/application.service ]; then
  echo "Stop application"
  sudo systemctl stop application
  sudo systemctl daemon-reload
fi
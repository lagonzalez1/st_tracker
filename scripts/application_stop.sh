#!/bin/bash
set -e
pwd

if [ -f "/home/systemd/system/application.service" ]; then
  echo "Stop application"
  sudo kill -9 $(sudo lsof -t -i:3333)
  sudo systemctl stop application
  sudo systemctl daemon-reload
fi
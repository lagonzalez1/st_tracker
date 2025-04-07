#!/bin/bash


pwd

sudo systemctl start application
sleep 2  # Allow service startup time
if ! systemctl is-active --quiet application; then
    echo "Error: Service failed to start"
    journalctl -u application -n 30 --no-pager
    exit 1
fi

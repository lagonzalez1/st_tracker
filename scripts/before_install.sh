#!/bin/bash

echo "Install secrets"
sudo yum update -y


sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt update
sudo apt install python3.11

sudo pwd

# Create env 
cd /home/ubuntu/app/scripts
sudo python3 deploy_env.py

#!/bin/bash

set -e # Fail and exit on error of any command

export WEBSCAN_VERSION=$(curl -s https://api.github.com/repos/thetillhoff/webscan/releases/latest | jq -r '.tag_name')
wget "https://github.com/thetillhoff/webscan/releases/download/${WEBSCAN_VERSION}/webscan_linux_amd64"
wget "https://github.com/thetillhoff/webscan/releases/download/${WEBSCAN_VERSION}/webscan_linux_amd64.sha256"
echo "$(cat webscan_linux_amd64.sha256) webscan_linux_amd64" | sha256sum --check --status
sudo install webscan_linux_amd64 /usr/local/bin/webscan # automatically sets rwxr-xr-x permissions
rm webscan_linux_amd64 webscan_linux_amd64.sha256

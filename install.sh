#!/bin/sh

set -e # Fail and exit on error of any command

REPO_OWNER='thetillhoff'
REPO_NAME='webscan'

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"

# Check if jq is available
if [ ! 'command -v jq' ]; then
  printf "jq is required to run this script"
  exit 1
fi

# Check if sha256sum is available
if [ ! 'command -v sha256sum' ]; then
  printf "sha256sum is required to run this script"
  exit 1
fi

# Check if curl or wget are available
if [ 'command -v curl' ]; then
  DOWNLOAD_FILE_CMD="curl -Lo"
  DOWNLOAD_BODY_CMD="curl -sL"
elif [ 'command -v wget' ]; then
  DOWNLOAD_FILE_CMD="wget -O"
  DOWNLOAD_BODY_CMD="wget -qO-"
else
  printf "Either curl or wget are required to run this script"
  exit 2
fi

LATEST_VERSION="$($DOWNLOAD_BODY_CMD https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest | jq -r '.tag_name')"
echo https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_VERSION}/webscan_${OS}_amd64
echo https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_VERSION}/webscan_${OS}_amd64.sha256
$DOWNLOAD_FILE_CMD ${REPO_NAME} "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_VERSION}/webscan_${OS}_amd64"
$DOWNLOAD_FILE_CMD ${REPO_NAME}.sha256 "https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_VERSION}/webscan_${OS}_amd64.sha256"
printf "$(cat ${REPO_NAME}.sha256) ${REPO_NAME}" | sha256sum --check --status
printf "Checksum validation complete, installing to /usr/local/bin/ ..."
sudo install ${REPO_NAME} /usr/local/bin/${REPO_NAME} # automatically sets rwxr-xr-x permissions
rm ${REPO_NAME} ${REPO_NAME}.sha256

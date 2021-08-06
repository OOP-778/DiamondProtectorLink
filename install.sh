#!/bin/bash

prompt=$(sudo -nv 2>&1)
if [ $? -eq 0 ]; then
  echo "Running this script as sudo!"
elif echo "$prompt" | grep -q '^sudo:'; then
    echo "This script must be ran with sudo privileges"
    exit 1;
else
  echo "This script must be ran with sudo privileges"
  exit 1;
fi

REPOSITORY="https://github.com/OOP-778/DiamondProtectorLink"
DIRECTORY="/usr/local/bin/DiamondProtectorLink"
MATERIAL_DIRECTORY=$DIRECTORY/material
OLD_DIR_EXISTS=0

if [ -d "$MATERIAL_DIRECTORY/" ]
then
    OLD_DIR_EXISTS=1
    echo "DiamondProtectorLink already found! Will be migrating..."
else
    echo "Base Directory does not exist, creating"
    sudo mkdir -p $MATERIAL_DIRECTORY/
fi;

cd $MATERIAL_DIRECTORY
if [ $OLD_DIR_EXISTS -eq 1 ]
then
  echo "Pulling changes from the repository"
  git pull
else
  git clone $REPOSITORY .
fi;

echo "Building DiamondProtectorLink"
if [ $OLD_DIR_EXISTS -eq 1 ]
then
  systemctl stop DiamondProtectorLink
fi;

go build

echo "Built successfully."
cp $MATERIAL_DIRECTORY/DiamondProtectorLink $DIRECTORY/DiamondProtectorLink
chmod u+x $DIRECTORY/DiamondProtectorLink

if [ $OLD_DIR_EXISTS -eq 0 ]
then
  echo "Creating default config"

  REDIS_HOSTNAME="localhost"
  REDIS_PORT="6379"
  REDIS_PASSWORD=""

  read -e -p "Redis Hostname: " -i "localhost" REDIS_HOSTNAME
  read -e -p "Redis port: " -i "6379" REDIS_PORT
  read -e -sp "Redis password: " -i "" REDIS_PASSWORD

  {
    echo "RedisHostname: $REDIS_HOSTNAME"
    echo "RedisPort: $REDIS_PORT"
    echo "RedisPassword: \"$REDIS_PASSWORD\""
  } >> $DIRECTORY/config.yml
fi;

if [ $OLD_DIR_EXISTS -eq 1 ]
then
  systemctl restart DiamondProtectorLink
else
  {
    echo "[Unit]"
    echo "Description=DiamondProtectorLink service"
    echo "After=syslog.target"
    echo " "
    echo "[Service]"
    echo "User=root"
    echo "WorkingDirectory=$DIRECTORY"
    echo "ExecStart=$DIRECTORY/DiamondProtectorLink"
    echo "Restart=on-failure"
    echo "StartLimitInterval=600"
    echo " "
    echo "[Install]"
    echo "WantedBy=multi-user.target"
  } >> /etc/systemd/system/DiamondProtectorLink.service
  systemctl enable --now DiamondProtectorLink
fi;
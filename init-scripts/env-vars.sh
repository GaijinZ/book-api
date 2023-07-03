#!/bin/sh

SOURCE_FILE="$HOME/.bash_profile"

touch "$SOURCE_FILE"

echo "USERS_SERVER_PORT=5000" > "$SOURCE_FILE"
echo "POSTGRES_USERS=postgres://tmosto:tmosto@172.19.0.2:5432/usersdb" >> "$SOURCE_FILE"
echo "POSTGRES_USER=tmosto" >> "$SOURCE_FILE"
echo "POSTGRES_PASSWORD=tmosto" >> "$SOURCE_FILE"
echo "POSTGRES_USERSDB=usersdb" >> "$SOURCE_FILE"
echo "POSTGRES_USERS_PORT=5432" >> "$SOURCE_FILE"
echo "POSTGRES_CONTAINER_IP=172.19.0.2" >> $SOURCE_FILE
echo "NETWORK=postgres" >> "$SOURCE_FILE"
echo "GOPATH=$HOME/go" >> "$SOURCE_FILE"
echo "PATH=$PATH:/usr/local/go/bin" >> "$SOURCE_FILE"
echo "SECRET_KEY=mysecretkeyshh" >> "$SOURCE_FILE"

echo "All variables has been saved in: $SOURCE_FILE"

source $SOURCE_FILE
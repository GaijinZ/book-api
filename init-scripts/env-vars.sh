#!/bin/sh

SOURCE_FILE="$HOME/.bash_profile"

sudo touch "$SOURCE_FILE"

chmod 777 "$SOURCE_FILE"

echo "USERS_SERVER_PORT=:5000" > "$SOURCE_FILE"
echo "POSTGRES_USERS=postgres://tmosto:tmosto@172.18.0.2:5432/usersdb" >> "$SOURCE_FILE"
echo "GOPATH=$HOME/go" >> "$SOURCE_FILE"
echo "PATH=$PATH:/usr/local/go/bin" >> "$SOURCE_FILE"
echo "SECRET_KEY=mysecretkeyshh" >> "$SOURCE_FILE"

echo "All variables has been saved in: $SOURCE_FILE"

source $SOURCE_FILE
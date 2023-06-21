#!/bin/bash
SOURCE_FILE="$HOME/.bash_profile"

touch "$SOURCE_FILE"

sudo chmod 777 "$SOURCE_FILE"

echo "export USERS_SERVER_PORT=:5000" > "$SOURCE_FILE"
echo "export POSTGRES_USERS=postgres://tmosto:tmosto@localhost:5432/usersdb" >> "$SOURCE_FILE"
echo "export GOPATH=$HOME/go" >> "$SOURCE_FILE"
echo "export PATH=$PATH:/usr/local/go/bin" >> "$SOURCE_FILE"

echo "All variables has been saved in: $SOURCE_FILE"

source $SOURCE_FILE
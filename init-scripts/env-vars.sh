#!/bin/sh

SOURCE_FILE="$HOME/.bash_profile"

touch "$SOURCE_FILE"

chmod +x init-scripts/env-vars.sh

echo "export USERS_SERVER_PORT=5000" > "$SOURCE_FILE"
echo "export POSTGRES_USERS=postgres://tmosto:tmosto@localhost:5432/usersdb" >> "$SOURCE_FILE"
echo "export POSTGRES_USER=tmosto" >> "$SOURCE_FILE"
echo "export POSTGRES_PASSWORD=tmosto" >> "$SOURCE_FILE"
echo "export POSTGRES_USERSDB=usersdb" >> "$SOURCE_FILE"
echo "export POSTGRES_USERS_PORT=5432" >> "$SOURCE_FILE"
echo "export GOPATH=$HOME/go" >> "$SOURCE_FILE"
echo "export PATH=$PATH:/usr/local/go/bin" >> "$SOURCE_FILE"
echo "export SECRET_KEY=mysecretkeyshh" >> "$SOURCE_FILE"

echo "All variables has been saved in: $SOURCE_FILE"

. $SOURCE_FILE
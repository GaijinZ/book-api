#!/bin/bash
SOURCE_FILE="/home/vagrant/.bash_profile"

touch "$SOURCE_FILE"

echo "USERS_SERVER_PORT=":5000"" >> "$SOURCE_FILE"
echo "POSTGRES_USERS="postgres://tmosto:tmosto@localhost:5432/usersdb"" >> "$SOURCE_FILE"
echo "GOPATH=$HOME/go" >> "$SOURCE_FILE"
echo "PATH=$PATH:/usr/local/go/bin" >> "$SOURCE_FILE"

# echo "All variables has been saved in: $SOURCE_FILE"

source $SOURCE_FILE
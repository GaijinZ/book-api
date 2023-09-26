#!/bin/sh

SOURCE_FILE="$HOME/.bash_profile"

touch "$HOME/.bash_profile"

echo "export GATEWAY_SERVER_PORT=8080" > "$SOURCE_FILE"
echo "export USERS_SERVER_PORT=5000" >> "$SOURCE_FILE"
echo "export BOOKS_SERVER_PORT=5001" >> "$SOURCE_FILE"
echo "export POSTGRES_USERS=postgres://tmosto:tmosto@localhost:5433/usersdb" >> "$SOURCE_FILE"
echo "export POSTGRES_BOOKS=postgres://tmosto:tmosto@localhost:5434/booksdb" >> "$SOURCE_FILE"
echo "export POSTGRES_USER=tmosto" >> "$SOURCE_FILE"
echo "export POSTGRES_PASSWORD=tmosto" >> "$SOURCE_FILE"
echo "export POSTGRES_USERSDB=usersdb" >> "$SOURCE_FILE"
echo "export POSTGRES_USERS_PORT=5433" >> "$SOURCE_FILE"
echo "export POSTGRES_BOOKSDB=booksdb" >> "$SOURCE_FILE"
echo "export POSTGRES_BOOKS_PORT=5434" >> "$SOURCE_FILE"
echo "export POSTGRES_USERS_CONTAINER_IP=172.19.0.2" >> "$SOURCE_FILE"
echo "export POSTGRES_BOOKS_CONTAINER_IP=172.19.0.3" >> "$SOURCE_FILE"
echo "export GOPATH=/usr/local/go" >> "$SOURCE_FILE"
echo "export GOROOT=/usr/local/go" >> "$SOURCE_FILE"
echo "export PATH=$PATH:/usr/local/go/bin" >> "$SOURCE_FILE"
echo "export SECRET_KEY=mysecretkeyshh" >> "$SOURCE_FILE"
echo "export NETWORK=bookapi_library" >> "$SOURCE_FILE"
echo "export GOPATH=$HOME/go" >> "$SOURCE_FILE"

echo "All variables has been saved in: $SOURCE_FILE"

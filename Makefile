include $(HOME)/.bash_profile

network:
	docker network create --subnet=172.19.0.0/16 library

build_users_postgres:
	docker build -t users-postgres -f internal/users/postgres/Dockerfile .

build_books_postgres:
	docker build -t books-postgres -f internal/books/postgres/Dockerfile .

run_users_postgres:
	docker run --network $(NETWORK) --name users_postgres --ip $(POSTGRES_USERS_CONTAINER_IP) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_USERSDB) -p 5433:5432 -d users-postgres

run_books_postgres:
	docker run --network $(NETWORK) --name books_postgres --ip $(POSTGRES_BOOKS_CONTAINER_IP) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_BOOKSDB) -p $(POSTGRES_BOOKS_PORT):5432 -d books-postgres

build_users:
	docker build --no-cache -t userapi-im -f internal/users/Dockerfile .

run_users:
	docker run --network $(NETWORK) --name users -p $(USERS_SERVER_PORT):$(USERS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh userapi-im

build_books:
	docker build --no-cache -t bookapi-im -f internal/books/Dockerfile .

run_books:
	docker run --network $(NETWORK) --name books -p $(BOOKS_SERVER_PORT):$(BOOKS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh bookapi

build_gateway:
	docker build --no-cache -t gateway -f internal/gateway/Dockerfile .

run_gateway:
	docker run --network $(NETWORK) --name gateway -p $(GATEWAY_SERVER_PORT):$(GATEWAY_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh gateway

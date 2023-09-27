include $(HOME)/.bash_profile

network:
	docker network create --subnet=172.19.0.0/16 library

build:
	docker build -t users-postgres -f users/postgres/Dockerfile .
	docker build -t books-postgres -f books/postgres/Dockerfile .
	docker build --no-cache -t userapi-im -f users/Dockerfile .
	docker build --no-cache -t bookapi-im -f books/Dockerfile .
	docker build --no-cache -t gateway -f gateway/Dockerfile .

run_db:
	docker run --network $(NETWORK) --name users_postgres --ip $(POSTGRES_USERS_CONTAINER_IP) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_USERSDB) -p $(POSTGRES_USERS_PORT):$(POSTGRES_USERS_PORT) -d users-postgres
	docker run --network $(NETWORK) --name books_postgres --ip $(POSTGRES_BOOKS_CONTAINER_IP) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_BOOKSDB) -p $(POSTGRES_BOOKS_PORT):$(POSTGRES_BOOKS_PORT) -d books-postgres

run_users:
	docker run --network $(NETWORK) --name users -p $(USERS_SERVER_PORT):$(USERS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh userapi-im

run_books:
	docker run --network $(NETWORK) --name books -p $(BOOKS_SERVER_PORT):$(BOOKS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh bookapi

run_gateway:
	docker run --network $(NETWORK) --name gateway -p $(GATEWAY_SERVER_PORT):$(GATEWAY_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh gateway

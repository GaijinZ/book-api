include $(HOME)/.bash_profile

network:
	docker network create --subnet=172.19.0.0/16 library

build:
	docker build --no-cache -t users-im -f users/cmd/Dockerfile .
	docker build --no-cache -t books-im -f books/cmd/Dockerfile .
	docker build --no-cache -t shops-im -f shops/cmd/Dockerfile .
	docker build --no-cache -t transactions-im -f transactions/cmd/Dockerfile .
	docker build --no-cache -t gateway-im -f gateway/cmd/Dockerfile .

run_psql:
	docker run -d -e POSTGRES_USER=$(BOOKAPI_POSTGRES_USER) -e POSTGRES_PASSWORD=$(BOOKAPI_POSTGRES_USER) -e POSTGRES_DB=$(BOOKAPI_POSTGRES_BOOKSDB) -p $(BOOKAPI_POSTGRES_BOOKS_PORT):$(BOOKAPI_POSTGRES_BOOKS_PORT) --name postgres postgres

run_rabbit:
	docker run -d --network $(BOOKAPI_NETWORK) --name rabbitmq rabbitmq:3

run_users:
	docker run --network $(BOOKAPI_NETWORK) --name users -p $(BOOKAPI_USERS_SERVER_PORT):$(BOOKAPI_USERS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh users-im

run_books:
	docker run --network $(BOOKAPI_NETWORK) --name books -p $(BOOKAPI_BOOKS_SERVER_PORT):$(BOOKAPI_BOOKS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh bookapi

run_shops:
	docker run --network $(BOOKAPI_NETWORK) --name gateway -p $(BOOKAPI_SHOPS_SERVER_PORT):$(BOOKAPI_SHOPS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh gateway

run_transactions:
	docker run --network $(BOOKAPI_NETWORK) --name gateway -p $(BOOKAPI_TRANSACTIONS_SERVER_PORT):$(BOOKAPI_TRANSACTIONS_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh gateway

run_gateway:
	docker run --network $(BOOKAPI_NETWORK) --name gateway -p $(BOOKAPI_GATEWAY_SERVER_PORT):$(BOOKAPI_GATEWAY_SERVER_PORT) --env-file init-scripts/env-vars-docker.sh gateway

counterfeiter:
	counterfeiter users/repository UsererRepository
	counterfeiter users/repository AutherRepository
	counterfeiter books/repository BookerRepository
	counterfeiter transactions/repository TransactionerRepository

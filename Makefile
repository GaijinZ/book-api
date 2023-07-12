include ${HOME}/.bash_profile

build_postgres:
	docker build --no-cache -t postgres -f database/postgres/Dockerfile .

run_postgres:
	docker run --network ${NETWORK} --name postgres --ip ${POSTGRES_CONTAINER_IP} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -e POSTGRES_DB=${POSTGRES_USERSDB} -p ${POSTGRES_USERS_PORT}:${POSTGRES_USERS_PORT} -d postgres

build_users:
	docker build --no-cache -t userapi -f internal/users/Dockerfile .

run_users:
	docker run --network ${NETWORK} --name users -p ${USERS_SERVER_PORT}:${USERS_SERVER_PORT} --env-file init-scripts/env-vars-docker.sh userapi

build_books:
	docker build --no-cache -t bookapi -f internal/books/Dockerfile .

run_books:
	docker run --network ${NETWORK} --name books -p ${BOOKS_SERVER_PORT}:${BOOKS_SERVER_PORT} --env-file init-scripts/env-vars-docker.sh bookapi
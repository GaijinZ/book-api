## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Requirements](#requirements)
* [Run app](#run)

## General info
Book API app is used to add, update, delete and view all books by user in Postgres database.
	
## Technologies
Project is created with:
* Golang 1.20
* Docker
* Postgres
* Redis
* RabbitMQ
* Logrus
* Swagger

## Run app
* Vagrant
* Virtual Box

## Run app
```
git clone https://github.com/GaijinZ/book-api.git

vagrant up

vagrant ssh
```

Build images
```
make build_postgres

make build_users

make build_books
```

Run containers
```
make run_postgres

make run_uers

make run_books
```

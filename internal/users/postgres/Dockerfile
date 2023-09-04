FROM postgres:latest

COPY init-scripts/users/docker-postgres-init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
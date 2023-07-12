FROM postgres:latest

COPY init-scripts/db/docker-postgres-init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
FROM postgres:15

COPY ./init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

FROM postgres:latest

ADD migrations/*.up.sql /docker-entrypoint-initdb.d
ADD seed.sql /docker-entrypoint-initdb.d

ARG POSTGRES_PASSWORD
ENV POSTGRES_PASSWORD ${POSTGRES_PASSWORD}
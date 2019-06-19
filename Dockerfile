FROM golang:latest
MAINTAINER wancat
WORKDIR /app
RUN apt-get update && apt-get install -y mariadb-client
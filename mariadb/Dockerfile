FROM mariadb:latest
MAINTAINER wancat
COPY init.sql /docker-entrypoint-initdb.d
WORKDIR /app
RUN mysqld --user=root --character-set-server=utf8 --collation-server=utf8_general_ci 
RUN echo "Asia/Taipei" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata
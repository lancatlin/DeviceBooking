FROM mariadb:latest
MAINTAINER wancat lancatlin@pm.me
ARG DB_NAME
ARG DB_USER
ARG DB_PASSWORD
ENV MYSQL_DATABASE ${DB_NAME}
ENV MYSQL_USER ${DB_USER}
ENV MYSQL_PASSWORD ${DB_PASSWORD}
ENV MYSQL_RANDOM_ROOT_PASSWORD "yes"
COPY ./init.sql /docker-entrypoint-initdb.d/
RUN echo "Asia/Shanghai" > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata

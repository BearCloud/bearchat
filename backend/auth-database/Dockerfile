FROM mysql/mysql-server:latest

ENV MYSQL_ROOT_PASSWORD=root

COPY initdb.sql /docker-entrypoint-initdb.d/

EXPOSE 3306

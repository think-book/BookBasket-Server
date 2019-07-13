FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD 'root'
ENV TZ 'Asia/Tokyo'

WORKDIR /database

COPY ./conf.d /etc/mysql/conf.d
COPY ./sql/* /docker-entrypoint-initdb.d/

EXPOSE 3306
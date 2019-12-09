FROM mysql:5.7

ENV MYSQL_ROOT_PASSWORD 'test'
ENV TZ 'Asia/Tokyo'

WORKDIR /database

COPY ./conf.d /etc/mysql/conf.d
COPY ./sql/* /docker-entrypoint-initdb.d/

EXPOSE 3306
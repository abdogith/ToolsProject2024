

FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD=abdomysql2001
ENV MYSQL_DATABASE=userdb

COPY userdb.sql /docker-entrypoint-initdb.d/
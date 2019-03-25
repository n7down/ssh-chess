FROM library/postgres

ENV POSTGRES_USER root 
ENV POSTGRES_PASSWORD password 
ENV POSTGRES_DB chess

COPY schema.sql /docker-entrypoint-initdb.d/1.sql
COPY data/players.sql /docker-entrypoint-initdb.d/2.sql
COPY data/games.sql /docker-entrypoint-initdb.d/3.sql

FROM golang:alpine3.10
lABEL maintainer="Pavel Sumin p.sumin7@gmail.com"

ENV TOKEN=""
ENV LOGIN=""
ENV PASSWORD=""
ENV NAME=""
ENV HOST=""
ENV PORT=""

COPY *.sql /docker-entrypoint-initdb.d/
COPY . /code
WORKDIR /code

EXPOSE 5432

RUN go build -o tg_bot

CMD API_TOKEN=${TOKEN} DB_HOST=${HOST} DB_PORT=${PORT} DB_LOGIN=${LOGIN} DB_PASSWORD=${PASSWORD} DB_NAME=${NAME} ./tg_bot
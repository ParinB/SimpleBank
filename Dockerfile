#Build stage
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.13
WORKDIR /app
COPY app.env .
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY db/migration ./migration
COPY start.sh .
COPY wait-for.sh .


EXPOSE 8080
#note cmd will be ignored  since in the compose file added an entry  point
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
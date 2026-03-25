FROM golang:1.20-alpine AS build
RUN apk add --no-cache git
WORKDIR /src
COPY . .
RUN go build -o /pgbadger ./...

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /pgbadger /pgbadger
EXPOSE 8080
ENTRYPOINT ["/pgbadger", "server"]

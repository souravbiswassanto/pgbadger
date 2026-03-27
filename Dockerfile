FROM golang:1.25-alpine AS build
RUN apk add --no-cache git
WORKDIR /src
COPY . .
RUN go build -o /pgbadger-server .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates perl pgbadger curl
COPY --from=build /pgbadger-server /pgbadger-server
EXPOSE 2385
ENTRYPOINT ["/pgbadger-server", "server"]

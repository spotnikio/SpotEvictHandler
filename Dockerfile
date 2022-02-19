FROM golang:1.17-alpine as dev-env

WORKDIR /app

FROM dev-env as build-env
COPY go.mod /go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 go build -o /main

FROM alpine:3.10 as runtime

COPY --from=build-env /main /usr/local/bin/main
RUN chmod +x /usr/local/bin/main

ENTRYPOINT ["main"]
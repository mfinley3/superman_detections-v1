FROM golang:1.14.2-buster as build

WORKDIR /go/src/

COPY go.mod .
COPY go.sum .
COPY cmd ./cmd
COPY internal ./internal

RUN go mod vendor
RUN go build ./cmd/api/main.go

FROM debian:buster 

WORKDIR /go

COPY --from=build go/src/main .
RUN mkdir resources
COPY resources/GeoLite2-City.mmdb ./resources

CMD ["./main"]

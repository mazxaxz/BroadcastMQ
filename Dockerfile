FROM golang:1.14 AS dependencies
WORKDIR /src
COPY go.mod .
COPY go.sum .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go mod download

FROM dependencies AS builder
WORKDIR /src
COPY . .
WORKDIR /src/cmd
RUN go build -ldflags="-w -s" -installsuffix bmq -tags=jsoniter -o /out/bmq .

FROM debian:stretch-slim
ENV DEBIAN_FRONTEND noninteractive
# add user for networking
WORKDIR /app
COPY --from=builder /out/app .

EXPOSE 8080
ENTRYPOINT ["/app/bmq"]
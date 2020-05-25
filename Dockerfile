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
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl net-tools \
    && apt-get clean -y \
    && apt-get autoremove -y \
    && rm -rf /tmp/* /var/tmp/* \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/bmq .

EXPOSE 8080
ENTRYPOINT ["/app/bmq"]
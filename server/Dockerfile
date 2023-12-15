# Build the application from source
FROM docker.io/library/golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod ./

COPY main.go ./
ADD logger /app/logger
ADD server_if /app/server_if
ADD tcp_server /app/tcp_server
ADD udp_server /app/udp_server

RUN ls -la /app

RUN GOOS=linux go build -o /ha-tcp-udp

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /ha-tcp-udp /ha-tcp-udp

EXPOSE 12000
EXPOSE 13000

USER nonroot:nonroot

ENTRYPOINT ["/ha-tcp-udp"]

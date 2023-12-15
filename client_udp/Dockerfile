# Build the application from source
FROM docker.io/library/golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod ./

COPY main.go ./

RUN GOOS=linux go build -o /udp-client

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /udp-client /udp-client

USER nonroot:nonroot

ENTRYPOINT ["/udp-client"]

# BUILD STAGE

FROM golang:1.23 AS build-stage

COPY . /app
WORKDIR /app

ARG VERSION=0.0.0-dev

RUN go mod download all \
    && CGO_ENABLED=0 GOOS=linux go build \
        -ldflags="-X main.version=${VERSION}" \
        -o ./opsrunner \
        main.go

# APP STAGE

FROM alpine:3.21 AS app-stage

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add ca-certificates tzdata \
    && update-ca-certificates

COPY --from=build-stage /app/opsrunner /opsrunner

ENTRYPOINT [ "/opsrunner" ]
CMD ["/opsrunner", "--help"]

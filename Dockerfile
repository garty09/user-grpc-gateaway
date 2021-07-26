FROM golang:alpine AS build
RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            gcc \
            musl-dev \
            ca-certificates && \
    rm -rf /var/cache/apk/*

# install migrate which will be used by entrypoint.sh to perform DB migration
ARG MIGRATE_VERSION=4.14.1
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./
RUN go mod download
RUN go mod verify

# copy source files and build the binary
COPY . .
RUN go build \
    -mod=readonly \
    -a -o /app/server ./cmd/user-test


FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
RUN mkdir -p /var/log/app
WORKDIR /app/
COPY --from=build /usr/local/bin/migrate /usr/local/bin
COPY --from=build /app/testdata ./migrations/
COPY --from=build /app/server .
COPY --from=build /app/config/*.yml ./config/
RUN ls -la
ENTRYPOINT ["/app/server"]
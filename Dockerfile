# build stage
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-bullseye AS build

WORKDIR /opt/build

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go test ./... \
    && go build -a -tags 'noplugins nossl netgo' -ldflags '-s -w' -o journey

# artefact stage
# hadolint ignore=DL3007
FROM gcr.io/distroless/base-debian11:latest

USER 1000
COPY --from=build --chown=1000:1000 /opt/build/journey  /usr/local/bin/journey
WORKDIR /opt/data
COPY --from=build --chown=1000:1000 /opt/build/built-in ./built-in
COPY --from=build --chown=1000:1000 /opt/build/content  ./content
CMD ["journey"]

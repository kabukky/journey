# build stage
FROM golang:1.19-bullseye AS build

WORKDIR /opt/build

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go test ./... \
    && go build -a -tags "noplugins nossl netgo linux" -ldflags '-s -w' -o journey

# artefact stage
FROM debian:bullseye-slim

COPY --from=build /opt/build/journey  /usr/local/bin/journey
USER nobody
WORKDIR /opt/data
COPY --from=build --chown=nobody:nobody /opt/build/built-in ./built-in
COPY --from=build --chown=nobody:nobody /opt/build/content  ./content
CMD ["journey"]

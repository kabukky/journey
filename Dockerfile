# build stage
FROM golang:1.18 AS build

WORKDIR /opt/app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go test ./... \
    && go build -a -tags "noplugins nossl netgo linux" -ldflags '-s -w' -o journey

# final stage
FROM debian:buster-slim

COPY --from=build /opt/app/journey  /usr/local/bin/journey
USER nobody
WORKDIR /opt/data
COPY --from=build --chown=nobody:nobody /opt/app/built-in ./built-in
COPY --from=build --chown=nobody:nobody /opt/app/content  ./content
CMD ["journey"]

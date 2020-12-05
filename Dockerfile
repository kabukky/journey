# build stage
FROM golang:1.14 AS build

WORKDIR /opt/journey
COPY . .
RUN git -c http.sslVerify=false submodule update --init --recursive
RUN go test ./...
RUN go build -a -tags "noplugins nossl netgo" -ldflags '-w' -o journey

# final stage
# hadolint ignore=DL3007
FROM ubuntu:18.04
WORKDIR /opt/journey
COPY --from=build /opt/journey/journey /opt/journey/
COPY --from=build /opt/journey/built-in /opt/journey/
COPY --from=build /opt/journey/config.yaml /opt/journey/
COPY --from=build /opt/journey/content /opt/journey/
CMD ["./journey"]

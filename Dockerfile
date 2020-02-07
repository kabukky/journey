# build stage
FROM golang:1.13 AS build

WORKDIR /opt/journey
COPY . .
RUN git -c http.sslVerify=false submodule update --init --recursive
RUN GOPROXY=http://172.17.0.1:8080 go build -a -o journey

# final stage
# hadolint ignore=DL3007
FROM ubuntu:18.04
WORKDIR /opt/journey
COPY --from=build /opt/journey/journey /opt/journey/
COPY --from=build /opt/journey/built-in /opt/journey/
COPY --from=build /opt/journey/config.json /opt/journey/
COPY --from=build /opt/journey/content /opt/journey/
CMD ["./journey"]

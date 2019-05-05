# build stage
FROM golang:1.12.4 AS build-env
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin
RUN go get github.com/kabukky/journey
RUN rm -rf /go/src/github.com/kabukky/journey
COPY . /go/src/github.com/kabukky/journey
WORKDIR /go/src/github.com/kabukky/journey
RUN git submodule update --init --recursive
RUN go get -d github.com/dimfeld/httptreemux
RUN go build

# final stage
# hadolint ignore=DL3007
FROM alpine:latest
WORKDIR /app
COPY --from=build-env /go/src/github.com/kabukky/journey/journey /app/
COPY --from=build-env /go/src/github.com/kabukky/journey/built-in /app/built-in
COPY --from=build-env /go/src/github.com/kabukky/journey/config.json /app/config.json
COPY --from=build-env /go/src/github.com/kabukky/journey/content /app/content
ENTRYPOINT ["./journey"]

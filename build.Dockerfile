# rm -f journey ; docker build -t tmp -f build.Dockerfile . && docker run -it --rm -v $(pwd):/mnt tmp cp journey /mnt/
# rm -f journey ; docker build --platform linux/amd64 -t tmp -f build.Dockerfile . && docker run --platform linux/amd64 -it --rm -v $(pwd):/mnt tmp cp journey /mnt/
FROM ubuntu:18.04

# hadolint ignore=DL3008
RUN apt-get update \
    && apt-get install --no-install-recommends -y \
        software-properties-common \
    && add-apt-repository ppa:longsleep/golang-backports \
    && apt-get install --no-install-recommends -y \
        git \
        golang-go \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/journey

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go test ./... \
    && go build -a -tags "noplugins nossl netgo" -ldflags '-s -w' -o journey
